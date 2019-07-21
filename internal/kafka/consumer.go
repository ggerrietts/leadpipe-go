package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Use ConsumerGroup from calling code

type Consumer struct {
	ready chan struct{}
}

// Setup is run at the beginning of a new session
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim begins a loop
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// probably write these to a channel, eh
		log.WithField("key", message.Key).WithField("value", message.Value).Debug("message received")
		session.MarkMessage(message, "")
	}
	return nil
}
