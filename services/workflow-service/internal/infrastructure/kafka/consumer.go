package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type WorkflowConsumer struct {
	reader *kafka.Reader
}

func NewWorkflowConsumer(broker, topic string) *WorkflowConsumer {
	return &WorkflowConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: "workflow-service",
		}),
	}
}

func (c *WorkflowConsumer) Start() {
	go func() {
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("kafka read error : ", err)
				continue
			}
			log.Println("Kafka message : ", string(msg.Value))

			var event map[string]interface{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("json error : ", err)
			}
		}
	}()
}
