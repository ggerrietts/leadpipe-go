package main

import (
	"fmt"
	"strings"

	"github.com/ggerrietts/leadpipe-go/internal/kafka"

	"github.com/ggerrietts/leadpipe-go/internal/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Log debug level or above.
	log.SetLevel(log.DebugLevel)

	// Log the function from which we call
	log.SetReportCaller(true)
}

func emitBootBanner(brokers []string, topic string) {
	msg := `
	 __    ____   __   ____  ____  __  ____  ____  _  ____  ____   __    ___  ____  ____  ____ 
	(  )  (  __) / _\ (    \(  _ \(  )(  _ \(  __)(_)(  _ \(  _ \ /  \  / __)(  __)/ ___)/ ___)
	/ (_/\ ) _) /    \ ) D ( ) __/ )(  ) __/ ) _)  _  ) __/ )   /(  O )( (__  ) _) \___ \\___ \
	\____/(____)\_/\_/(____/(__)  (__)(__)  (____)(_)(__)  (__\_) \__/  \___)(____)(____/(____/

	[-] Brokers at: %v
	[-] Logging level: %v
	[-] Topic: %v
	`
	fmt.Printf(msg, brokers, log.GetLevel().String(), topic)
}

func main() {
	config.Load()
	cfg := config.Load()
	brokers := strings.Split(cfg.GetString(config.KafkaBrokers), ",")
	topic := cfg.GetString(config.KafkaTopic)

	cgConfig := kafka.ConsumerGroupConfig{
		Brokers:       brokers,
		Version:       cfg.GetString(config.KafkaClusterVersion),
		Topic:         topic,
		ConsumerGroup: cfg.GetString(config.KafkaConsumerGroup),
	}

	consumerGroup := kafka.NewConsumerGroup(cgConfig)

	emitBootBanner(brokers, topic)
	consumerGroup.Consume()
}
