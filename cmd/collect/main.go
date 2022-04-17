package main

import (
	"github.com/ggerrietts/leadpipe-go/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/pkg/collect"
	"github.com/ggerrietts/leadpipe-go/pkg/kafka"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	config.Load()
	cfg := config.Load()
	kafkaAddr := cfg.GetString(config.KafkaAddr)
	listenAddr := cfg.GetString(config.CollectorAddr)
	topic := cfg.GetString(config.ProcessTopic)

	logrus.WithFields(logrus.Fields{
		"listenAddr": listenAddr,
		"kafkaAddr":  kafkaAddr,
		"topic":      topic,
	}).Info("starting")

	producer := kafka.NewProducer(kafkaAddr, topic)
	defer producer.Close()
	server := collect.NewServer(producer, listenAddr)

	server.Serve()
}
