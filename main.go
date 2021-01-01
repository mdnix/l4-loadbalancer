package main

import (
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mdnix/roundrobin"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	configFile string
)

func copyConnection(src, dst net.Conn, done chan bool) {
	_, err := io.Copy(src, dst)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("copy")
	}
	done <- true
}

func handleConn(src net.Conn, backend string) {
	defer src.Close()

	dst, err := net.Dial("tcp", backend)
	defer dst.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("unable to connect to to backend")
	}

	done := make(chan bool, 2)
	go copyConnection(src, dst, done)
	go copyConnection(dst, src, done)
	<-done
	<-done
}

func getNextBackendAddress(service *roundrobin.Service) (string, error) {
	nextBackend := service.Next()
	alive := isAlive(nextBackend)
	var count int
	for !alive {
		if count == len(service.Backends) {
			return "", errors.New("no healthy backends available")
		}
		nextBackend = service.Next()
		alive = isAlive(nextBackend)
		count++
	}
	return nextBackend.Address, nil
}

func init() {
	flag.StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.services)")
	initConfig()
}

func serve(bind string, backends []string) {
	service, err := roundrobin.NewService(backends)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to create service")
	}

	ln, err := net.Listen("tcp", bind)
	defer ln.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("no listener")
	}

	log.WithFields(log.Fields{
		"address": bind,
	}).Info("listener started")

	go scheduleHealthCheck(service, time.Minute*5)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("unable to accept connection")
			continue
		}

		go func() {
			if nextAddress, err := getNextBackendAddress(service); err == nil {
				writeBalancingFlow(conn.RemoteAddr().String(), conn.LocalAddr().String(), nextAddress)
				handleConn(conn, nextAddress)
			} else {
				conn.Close()
				log.WithFields(log.Fields{
					"error": err,
				}).Warn("unable to handle request")
			}
		}()
	}
}

func main() {
	for _, s := range config.Services {
		go serve(s.Bind, s.Backends)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	log.Info("shutting down loadbalancer")
}
