package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// ConsumerGroup is a single Kafka consumer
type ConsumerGroup struct {
	client sarama.ConsumerGroup
	topic  string
}

// ConsumerGroupConfig collects the options necessary to start a consumer group
type ConsumerGroupConfig struct {
	Brokers       []string
	Version       string
	Topic         string
	ConsumerGroup string
}

// NewConsumerGroup creates a Consumer given the broker addresses
func NewConsumerGroup(in ConsumerGroupConfig) *ConsumerGroup {

	version, err := sarama.ParseKafkaVersion(in.Version)
	if err != nil {
		log.WithError(err).Fatal("error parsing kafka version")
	}
	config := sarama.NewConfig()
	config.Version = version

	client, err := sarama.NewConsumerGroup(in.Brokers, in.ConsumerGroup, config)
	if err != nil {
		log.WithError(err).WithField("brokers", in.Brokers).WithField("cg", in.ConsumerGroup).Fatal("error starting consumer group")
	}

	return &ConsumerGroup{
		client: client,
		topic:  in.Topic,
	}
}

// Consume starts up the ConsumerGroup
func (cg *ConsumerGroup) Consume() {
	log.Debug("starting Consume")

	consumer := Consumer{
		ready: make(chan bool, 0),
	}
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			log.Debug("mmmm so hungry")
			topics := []string{cg.topic}
			if err := cg.client.Consume(ctx, topics, &consumer); err != nil {
				log.WithError(err).Panic("error from consumer")
			}

			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool, 0)
		}
	}()
	log.Debug("waiting for godot")
	<-consumer.ready
	log.Debug("consumer up and running")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.WithError(ctx.Err()).Warn("terminating: context canceled")
	case <-sigterm:
		log.Warn("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := cg.client.Close(); err != nil {
		log.WithError(err).Panic("Error closing client")
	}
}
