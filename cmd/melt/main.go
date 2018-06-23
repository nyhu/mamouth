package main

import (
	"fmt"
	"time"

	"github.com/nyhu/mamouth/c14"
	"github.com/nyhu/mamouth/entity"
	"github.com/nyhu/mamouth/kafka"
)

func main() {
	kafka, err := kafka.NewKafka("51.15.231.63:29092")
	if err != nil {
		fmt.Println("error: ", err)
	}

	fmt.Println("Starting crying Melt")

	batchChan := make(chan entity.KafkaBatch)
	go kafka.Produce("test-2", batchChan)
	topic := "test-2"
	startTime := time.Now()
	endTime := time.Now()
	safeId, err := c14.GetSafe(topic)

	c14.Melt(safeId, startTime, endTime, batchChan)
}
