package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	topic  string
	broker string
	writer *kafka.Writer
}

// สร้าง producer
func NewProducer(broker string, topic string) *Producer {
	producer := &Producer{
		topic:  topic,
		broker: broker,
	}

	if err := producer.ensureTopicExists(); err != nil {
		log.Println("kafka topic check failed:", err)
	}

	producer.writer = &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	log.Println("kafka producer initialized for topic:", topic)
	return producer
}

func (p *Producer) ensureTopicExists() error {
	conn, err := kafka.Dial("tcp", p.broker)
	if err != nil {
		log.Fatalf("failed to connect to kafka broker: %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		log.Fatalf("failed to read partitions: %v", err)
	}

	for _, partition := range partitions {
		if partition.Topic == p.topic {
			return nil
		}
	}
	// create topic
	if err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             p.topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		return fmt.Errorf("create topic failed: %w", err)
	}

	log.Println("kafka topic created : ", p.topic)
	return nil
}

// push message
func (p *Producer) Publish(key string, value interface{}) error {
	if p == nil || p.writer == nil {
		return fmt.Errorf("kafka producer not initialized")
	}
	messageBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
			Time:  time.Now(),
		})
}
