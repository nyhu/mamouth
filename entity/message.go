package entity

type KafkaMessage struct {
	Content []byte 	`json:"content"`
	Offset int64	`json:"offset"`
}

func NewKafkaMessage(content []byte, offset int64) KafkaMessage {
	return KafkaMessage{
		Content:content,
		Offset:offset,
	}
}

type KafkaBatch struct {
	Batch []KafkaMessage `json:"batch"`
}
