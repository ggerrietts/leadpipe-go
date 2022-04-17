package main

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

type Payload struct {
	Time   time.Time
	Sonnet string
	Weight int32
}

func RandomPayload() *Payload {
	return &Payload{
		Time:   time.Now(),
		Sonnet: LinePlease(),
		Weight: rand.Int31(),
	}
}

func main() {
	url := "http://localhost:8080/"
	for {
		payload := RandomPayload()
		foo, err := json.Marshal(payload)
		if err != nil {
			logrus.WithError(err).Fatal("could not even")
		}
		resp, err := http.Post(url, "application/json", bytes.NewReader(foo))
		if err != nil {
			logrus.WithError(err).Fatal("oopadaisy")
		}
		if resp.StatusCode != http.StatusOK {
			logrus.WithField("StatusCode", resp.StatusCode).Fatal("bad response")
		}
		logrus.WithField("payload", payload)
		time.Sleep(time.Second * 2)
	}
}
