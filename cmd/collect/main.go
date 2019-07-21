package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/collect"
	"github.com/ggerrietts/leadpipe-go/internal/kafka"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	producer := kafka.NewKafkaProducer([]string{"192.168.1.10:9092"})
	server := collect.NewServer(producer)
	log.Debug("preparing to serve")
	server.Serve()
}
