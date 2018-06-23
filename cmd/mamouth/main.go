package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nyhu/mamouth/c14"
	"github.com/nyhu/mamouth/kafka"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	relayChan := make(chan string, 3600)

	kafka, err := kafka.NewKafka("51.15.231.63:29092")
	if err != nil {
		fmt.Println("error: ", err)
	}
	topicName := "test-2"
	fmt.Println("starting fucking mamouth")
	go kafka.Consume(topicName, relayChan, signals)
	c14.Relay(topicName, relayChan, signals)
}
