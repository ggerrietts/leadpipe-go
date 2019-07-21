package config

import "github.com/spf13/viper"

// SetDefaults sets up default values for config
func SetDefaults() {
	viper.SetDefault(KafkaBrokers, "192.168.1.10:9092")
	viper.SetDefault(KafkaClusterVersion, "3.0.1")
	viper.SetDefault(KafkaConsumerGroup, "leadpipe")
	viper.SetDefault(KafkaTopic, "hits")
	viper.SetDefault(CollectorAddr, ":8080")
}
