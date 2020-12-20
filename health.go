package main

import (
	"net"
	"time"

	"github.com/mdnix/roundrobin"
	log "github.com/sirupsen/logrus"
)

func scheduleHealthCheck(s *roundrobin.Service, t time.Duration) {
	log.WithFields(log.Fields{
		"interval": t,
	}).Info("healthcheck started")
	healthCheck(s)
	for {
		<-time.After(t)
		go healthCheck(s)
	}
}

func healthCheck(s *roundrobin.Service) {
	for _, backend := range s.Backends {
		status := "up"
		alive := isBackendAlive(backend)
		setAlive(alive, backend)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]", backend.Address, status)
	}
}

func setAlive(alive bool, backend *roundrobin.Backend) {
	backend.Mu.Lock()
	backend.IsAlive = alive
	backend.Mu.Unlock()
}

func isAlive(backend *roundrobin.Backend) (alive bool) {
	backend.Mu.RLock()
	alive = backend.IsAlive
	backend.Mu.RUnlock()
	return alive
}

func isBackendAlive(backend *roundrobin.Backend) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", backend.Address, timeout)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warnf("unable to reach %s", backend.Address)
		return false
	}
	_ = conn.Close()
	return true
}
