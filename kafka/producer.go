package kafka

import (
	"github.com/nyhu/mamouth/entity"
	"github.com/Shopify/sarama"
	"time"
	"fmt"
)

func (k *Kafka) Produce(topicName string, batchChan chan entity.KafkaBatch) error {
	producer, err := k.newProducer()
	if err != nil {
		return err
	}

	for {
		batch := <- batchChan
		before := time.Now()

		var messages []*sarama.ProducerMessage
		for i, msg := range batch.Batch {
			messages = append(messages, prepareMessage(topicName, string(msg.Content)))
			if i % 50000 == 0 {
				fmt.Println("sending index", i)
				producer.SendMessages(messages)
				messages = []*sarama.ProducerMessage{}
			}
		}
		producer.SendMessages(messages)
		after := time.Now()
		fmt.Println("produced", len(batch.Batch), "message in", after.Unix() - after.Unix(), "second",
			after.UnixNano() - before.UnixNano(), "nanosecond")
	}
}

func (k *Kafka) newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{k.bootstrapServer}, config)

	return producer, err
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

