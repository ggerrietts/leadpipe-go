package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ggerrietts/leadpipe-go/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/kafka"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Log debug level or above.
	log.SetLevel(log.DebugLevel)

	// Log the function from which we call
	log.SetReportCaller(true)
}

func emitBootBanner(brokers []string, topic string, cGroup string) {
	msg := `
	 __    ____   __   ____  ____  __  ____  ____  _  ____  ____   __    ___  ____  ____  ____ 
	(  )  (  __) / _\ (    \(  _ \(  )(  _ \(  __)(_)(  _ \(  _ \ /  \  / __)(  __)/ ___)/ ___)
	/ (_/\ ) _) /    \ ) D ( ) __/ )(  ) __/ ) _)  _  ) __/ )   /(  O )( (__  ) _) \___ \\___ \
	\____/(____)\_/\_/(____/(__)  (__)(__)  (____)(_)(__)  (__\_) \__/  \___)(____)(____/(____/

	[-] Brokers at: %v
	[-] Logging level: %v
	[-] Topic: %v
	[-] Consumer Group: %v
`
	fmt.Printf(msg, brokers, log.GetLevel().String(), topic, cGroup)
}

func main() {
	config.Load()
	cfg := config.Load()
	brokers := strings.Split(cfg.GetString(config.KafkaBrokers), ",")
	topic := cfg.GetString(config.KafkaTopic)
	cGroup := cfg.GetString(config.KafkaConsumerGroup)
	buffDepth := cfg.GetInt(config.MessageChannelDepth)

	emitBootBanner(brokers, topic, cGroup)

	msgChan := make(chan *kafka.Message, buffDepth)

	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		ConsumerGroupID: cGroup,
		Topic:           topic,
		Brokers:         brokers,
		MessageChan:     msgChan,
	})
	defer consumer.Close()

	go func() {
		consumer.Consume(context.Background())
	}()

	select {}
}
