package main

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
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

// func main() {
// 	config.Load()
// 	cfg := config.Load()
// 	brokers := strings.Split(cfg.GetString(config.KafkaBrokers), ",")
// 	topic := cfg.GetString(config.KafkaTopic)

// 	emitBootBanner(brokers, topic)

// 	cgConfig := kafka.ConsumerGroupConfig{
// 		Brokers:       brokers,
// 		Version:       cfg.GetString(config.KafkaClusterVersion),
// 		Topic:         topic,
// 		ConsumerGroup: cfg.GetString(config.KafkaConsumerGroup),
// 	}

// 	consumerGroup := kafka.NewConsumerGroup(cgConfig)

// 	consumerGroup.Consume()
// }

var (
	// kafka
	kafkaBrokerUrl     string
	kafkaVerbose       bool
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string
)

func main() {
	kafkaBrokerUrl := "192.168.1.10:9092"
	// kafkaVerbose := true
	kafkaTopic := "hits"
	// kafkaConsumerGroup := "leadpipe"
	// kafkaClientId := "reader"

	brokers := strings.Split(kafkaBrokerUrl, ",")

	emitBootBanner(brokers, kafkaTopic)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		// GroupID: "leadpipe",
		Partition: 0,
		Topic:     kafkaTopic,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/offset %v/%v: %s = %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))
	}

}
