package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyhu/mamouth/entity"
	"github.com/nyhu/mamouth/kafka"
)

type caseMsg struct {
	SensorName string    `json:"sensorName"`
	Time       time.Time `json:"time"`
	Value      int       `json:"value"`
}

func makeMsg(i int) entity.KafkaMessage {
	encodedMsg, err := json.Marshal(caseMsg{
		SensorName: "[gossip]moumouth/decibel",
		Time:       time.Now(),
		Value:      i,
	})
	if err != nil {
		return entity.KafkaMessage{
			Offset:  int64(i),
			Content: []byte{'g', 'o', 's', 's', 'i', 'p'},
		}
	}
	return entity.KafkaMessage{
		Offset:  int64(i),
		Content: encodedMsg,
	}
}

func makeBatch() entity.KafkaBatch {
	batch := entity.KafkaBatch{}
	for i := 0; i < 500001; i++ {
		batch.Batch = append(batch.Batch, makeMsg(i))
	}
	return batch
}

func sendInLoop(bc chan entity.KafkaBatch) {
	for {
		fmt.Println("Sending 500000 messages")
		bc <- makeBatch()
	}
}

func main() {
	kafka, err := kafka.NewKafka("51.15.231.63:29092")
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("starting crying gossip")

	batchChan := make(chan entity.KafkaBatch)
	go sendInLoop(batchChan)
	if err := kafka.Produce("test-2", batchChan); err != nil {
		fmt.Println(err)
	}
}
