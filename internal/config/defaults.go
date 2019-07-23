package config

import (
	"runtime"

	"github.com/spf13/viper"
)

// SetDefaults sets up default values for config
func SetDefaults() {
	viper.SetDefault(KafkaBrokers, "192.168.1.10:9092")
	viper.SetDefault(KafkaClusterVersion, "5.3.0")
	viper.SetDefault(KafkaConsumerGroup, "leadpipe")
	viper.SetDefault(KafkaTopic, "hits")
	viper.SetDefault(CollectorAddr, ":8080")
	viper.SetDefault(MessageChannelDepth, 2*runtime.GOMAXPROCS(0))
}
