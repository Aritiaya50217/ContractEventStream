package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Start() {
	go func() {
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("error reading kafka message:", err)
				continue
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Println("error unmarshalling message:", err)
				continue
			}

			log.Printf("received message: key=%s value=%v", string(msg.Key), data)
		}
	}()
}
