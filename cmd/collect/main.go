package main

import (
	"fmt"
	"strings"

	"github.com/ggerrietts/leadpipe-go/internal/config"

	log "github.com/sirupsen/logrus"

	"github.com/ggerrietts/leadpipe-go/internal/collect"
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

func emitBootBanner(addr string, brokers []string, topic string) {
	msg := `
	 __    ____   __   ____  ____  __  ____  ____  _  ___  __   __    __    ____  ___  ____ 
	(  )  (  __) / _\ (    \(  _ \(  )(  _ \(  __)(_)/ __)/  \ (  )  (  )  (  __)/ __)(_  _)
	/ (_/\ ) _) /    \ ) D ( ) __/ )(  ) __/ ) _)  _( (__(  O )/ (_/\/ (_/\ ) _)( (__   )(  
	\____/(____)\_/\_/(____/(__)  (__)(__)  (____)(_)\___)\__/ \____/\____/(____)\___) (__) 

	[-] Booting on: %v
	[-] Brokers at: %v
	[-] Logging level: %v
	[-] Topic: %v
`
	fmt.Printf(msg, addr, brokers, log.GetLevel().String(), topic)
}

func main() {
	config.Load()
	cfg := config.Load()
	brokers := strings.Split(cfg.GetString(config.KafkaBrokers), ",")
	listenAddr := cfg.GetString(config.CollectorAddr)
	topic := cfg.GetString(config.ProcessTopic)

	emitBootBanner(listenAddr, brokers, topic)

	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()
	server := collect.NewServer(producer, listenAddr)

	server.Serve()
}
