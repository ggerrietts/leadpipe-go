package kafka

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/pb"
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
func (p *Producer) Send(ctx context.Context, hit *pb.Hit) {
	payload, err := proto.Marshal(hit)
	if err != nil {
		log.WithError(err).WithField("hit", hit).Error("failed to marshal hit")
		return
	}
	msg := kafka.Message{
		Key:   []byte(hit.Token),
		Value: payload,
	}
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.WithError(err).Error("failed to produce hit")
		return
	}
	log.WithField("hit", hit).Debug("produced hit")
}

// Close closes the writer
func (p *Producer) Close() {
	p.writer.Close()
}
