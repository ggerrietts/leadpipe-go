package kafka

import (
	"context"

	"github.com/ggerrietts/leadpipe-go/pkg/pb"
	"github.com/golang/protobuf/proto"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// Message models a key/value message
type Message struct {
	VisitorID string
	Payload   string
}

// Consumer coordinates the consuming of the queue
type Consumer struct {
	msgChan chan *pb.LeadpipeEvent
	reader  *kafka.Reader
}

// ConsumerConfig configures a Consumer
type ConsumerConfig struct {
	ConsumerGroupID string
	Brokers         []string
	Topic           string
	MessageChan     chan *pb.LeadpipeEvent
}

// NewConsumer creates a new consumer
func NewConsumer(cc ConsumerConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cc.Brokers,
		Topic:    cc.Topic,
		GroupID:  cc.ConsumerGroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{
		msgChan: cc.MessageChan,
		reader:  reader,
	}
}

// Close closes the reader
func (c *Consumer) Close() {
	c.reader.Close()
}

// MsgChan returns a message channel for consuming
func (c *Consumer) MsgChan() <-chan *pb.LeadpipeEvent {
	return c.msgChan
}

// Consume Mass Quantities
func (c *Consumer) Consume(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.WithError(err).Error("error reading message")
			c.reader.CommitMessages(ctx, m)
			continue
		}

		var msg pb.LeadpipeEvent
		err = proto.Unmarshal(m.Value, &msg)
		if err != nil {
			log.WithError(err).WithField("m.Value", m.Value).Error("error decoding message")
			c.reader.CommitMessages(ctx, m)
			continue
		}

		select {
		case <-ctx.Done():
			log.WithField("msg", msg).Info("context timeout")
			return
		// case c.msgChan <- &msg:
		default:
			log.WithField("msg", msg).Debug("message received loud and clear loud and clear")
		}
		c.reader.CommitMessages(ctx, m)
	}
}
