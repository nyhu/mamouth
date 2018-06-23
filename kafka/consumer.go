package kafka

import (
	"github.com/Shopify/sarama"
	"os"
	"fmt"
	"github.com/nyhu/mamouth/entity"
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	BASEPATH = "/tmp/"
)

type Kafka struct {
	bootstrapServer string
}

func NewKafka(bootstrapServer string) (*Kafka, error) {
	return &Kafka{
		bootstrapServer: bootstrapServer,
	}, nil
}

func (k *Kafka) Consume(topicName string, pathRelayChan chan string, signal chan os.Signal) (error) {
	err := ensureTopicDirectory(BASEPATH + topicName)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{k.bootstrapServer}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topicName, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	// Count how many entity processed

	// Get signal for finish
	msgCount := 0
	batch := entity.KafkaBatch{}
	doneCh := make(chan struct{})
	go
	func() {
		now := time.Now().Unix()
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)

			case msg := <-consumer.Messages():
				msgCount++
				batch.Batch = append(batch.Batch, entity.NewKafkaMessage(msg.Value, msg.Offset))
				// If a second is elapsed
				if time.Now().Unix() > now {
					fmt.Println("Sending messages", msgCount)
					msgCount = 0

					go dumpMessages(batch, topicName, now, pathRelayChan)
					batch = entity.KafkaBatch{}

					now = time.Now().Unix()
				}

			case <-signal:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	return nil
}

func ensureTopicDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func dumpMessages(messages entity.KafkaBatch, topicName string, now int64, pathRelayChan chan string) error {
	jsonM, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s%s/%v", BASEPATH, topicName, now)

	err = ioutil.WriteFile(path, jsonM, 0644)
	if err != nil {
		return err
	}

	pathRelayChan <- path
	return nil
}
