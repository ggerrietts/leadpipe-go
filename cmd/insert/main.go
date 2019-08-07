package main

import (
	"fmt"
	"strings"

	"github.com/ggerrietts/leadpipe-go/internal/config"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Log debug level or above.
	logrus.SetLevel(logrus.DebugLevel)

	// Log the function from which we call
	logrus.SetReportCaller(true)
}

func emitBootBanner(brokers []string, topic string, cGroup string, cHouse string) {
	msg := `
	 __    ____   __   ____  ____  __  ____  ____  _  __  __ _  ____  ____  ____  ____ 
	(  )  (  __) / _\ (    \(  _ \(  )(  _ \(  __)(_)(  )(  ( \/ ___)(  __)(  _ \(_  _)
	/ (_/\ ) _) /    \ ) D ( ) __/ )(  ) __/ ) _)  _  )( /    /\___ \ ) _)  )   /  )(  
	\____/(____)\_/\_/(____/(__)  (__)(__)  (____)(_)(__)\_)__)(____/(____)(__\_) (__) 

	[-] Brokers at: %v
	[-] Logging level: %v
	[-] Consume Topic: %v
	[-] Consumer Group: %v
	[-] Clickhouse DSN: %v
`
	fmt.Printf(msg, brokers, logrus.GetLevel().String(), topic, cGroup, cHouse)
}

func main() {
	cfg := config.Load()
	brokers := strings.Split(cfg.GetString(config.KafkaBrokers), ",")
	topic := cfg.GetString(config.InsertTopic)
	cGroup := cfg.GetString(config.InsertConsumerGroup)
	cHouse := cfg.GetString(config.ClickhouseAddr)
	emitBootBanner(brokers, topic, cGroup, cHouse)
}
