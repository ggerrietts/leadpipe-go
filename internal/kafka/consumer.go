package kafka

import (
	"context"

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
	msgChan chan *Message
	reader  *kafka.Reader
}

// ConsumerConfig configures a Consumer
type ConsumerConfig struct {
	ConsumerGroupID string
	Brokers         []string
	Topic           string
	MessageChan     chan *Message
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
func (c *Consumer) MsgChan() <-chan *Message {
	return c.msgChan
}

// Consume Mass Quantities
func (c *Consumer) Consume(ctx context.Context) {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.WithError(err).Error("error reading message")
			return
		}
		msg := Message{
			VisitorID: string(m.Key),
			Payload:   string(m.Value),
		}
		select {
		case <-ctx.Done():
			log.WithField("msg", msg).Info("context timeout")
			return
		case c.msgChan <- &msg:
			log.WithField("msg", msg).Debug("message received loud and clear loud and clear")
		}
		c.reader.CommitMessages(ctx, m)
	}
}
