package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Producer wraps the Kafka producer
type Producer interface {
	Send(string, string)
}

// KafkaProducer is a thin wrapper around a Sarama producer. It hides the Sarama API
type KafkaProducer struct {
	producer sarama.SyncProducer
}

// NewKafkaProducer creates a new producer given a list of broker addresses
func NewKafkaProducer(brokers []string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // wait for all in-sync replicas to ack
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	// TLS would go here

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.WithError(err).Fatal("failed to start sarama producer")
	}
	return &KafkaProducer{
		producer: producer,
	}
}

// Send wraps the sarama producer SendMessage
func (p *KafkaProducer) Send(visitorID string, payload string) {
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
