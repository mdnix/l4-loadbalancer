package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	tab       = "\t"
	newline   = "\n"
	separator = "------------------------\n"
)

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(a ...interface{}) {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprint(ew.w, a...)
}

func writeBalancingFlow(src, lb, dst string) error {
	switch config.Global.Output {
	case "dict":
		ew := &errWriter{w: os.Stdout}

		ew.write(
			"   TIMESTAMP: ", time.Now().Format(time.RFC3339), newline,
			"      SOURCE: ", src, newline,
			"LOADBALANCER: ", lb, newline,
			" DESTINATION: ", dst, newline,
			separator,
		)
		if ew.err != nil {
			return fmt.Errorf("failed to write out packet: %v", ew.err)
		}

	case "compact":
		_, err := fmt.Fprintf(os.Stdout,
			"%s: %s -> %s -> %s \n",
			time.Now().Format(time.RFC3339),
			src,
			lb,
			dst,
		)
		if err != nil {
			return fmt.Errorf("failed to write out packet: %v", err)
		}
	}
	return nil
}
