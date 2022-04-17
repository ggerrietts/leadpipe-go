package config

import (
	"runtime"

	"github.com/spf13/viper"
)

// SetDefaults sets up default values for config
func SetDefaults() {
	viper.SetDefault(KafkaAddr, "localhost:9092")
	viper.SetDefault(ProcessConsumerGroup, "leadpipe-hits")
	viper.SetDefault(ProcessTopic, "hits")
	viper.SetDefault(CollectorAddr, ":8080")
	viper.SetDefault(MessageChannelDepth, 2*runtime.GOMAXPROCS(0))
	viper.SetDefault(InsertConsumerGroup, "leadpipe-inserts")
	viper.SetDefault(InsertTopic, "inserts")
	viper.SetDefault(ClickhouseAddr, "localhost:8123")
}
