package kafka

import (
	"context"
	"github.com/golang/protobuf/proto"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/pkg/pb"
)

// Producer is a thin wrapper around a Sarama producer. It hides the Sarama API
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// NewProducer creates a new producer given a list of broker addresses
func NewProducer(broker string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

// Send wraps the sarama producer SendMessage
func (p *Producer) Send(ctx context.Context, evt *pb.LeadpipeEvent) {
	payload, err := proto.Marshal(evt)
	if err != nil {
		logrus.WithError(err).WithField("hit", evt).Error("failed to marshal event")
		return
	}
	msg := kafka.Message{
		Value: payload,
	}
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		logrus.WithError(err).Error("failed to produce event")
		return
	}
	logrus.WithField("event", evt).Debug("produced event")
}

// Close closes the writer
func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		logrus.WithError(err).Fatal("failed to close writer")
	}
}
