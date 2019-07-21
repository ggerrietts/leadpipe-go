package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Producer is a thin wrapper around a Sarama producer. It hides the Sarama API
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new producer given a list of broker addresses
func NewProducer(brokers []string, topic string) *Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // wait for all in-sync replicas to ack
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	// TLS would go here

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.WithError(err).Fatal("failed to start sarama producer")
	}
	return &Producer{
		producer: producer,
		topic:    topic,
	}
}

// Send wraps the sarama producer SendMessage
func (p *Producer) Send(visitorID string, payload string) {
	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "hits",
		Key:   sarama.StringEncoder(visitorID),
		Value: sarama.StringEncoder(payload),
	})
	if err != nil {
		log.WithError(err).Error("failed to produce event")
		return
	}
	log.WithFields(log.Fields{
		"partition": partition,
		"offset":    offset,
	}).Debug("produced event")
}
