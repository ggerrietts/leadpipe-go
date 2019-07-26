package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// Producer is a thin wrapper around a Sarama producer. It hides the Sarama API
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// NewProducer creates a new producer given a list of broker addresses
func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

// Send wraps the sarama producer SendMessage
func (p *Producer) Send(ctx context.Context, visitorID string, payload string) {
	msg := kafka.Message{
		Key:   []byte(visitorID),
		Value: []byte(payload),
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.WithError(err).Error("failed to produce event")
		return
	}
	log.WithField("visitorID", visitorID).Debug("produced event")
}

// Close closes the writer
func (p *Producer) Close() {
	p.writer.Close()
}
