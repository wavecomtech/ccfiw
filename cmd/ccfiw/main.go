package main

import (
	"ccfiw/cmd/ccfiw/cmd"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

}

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
