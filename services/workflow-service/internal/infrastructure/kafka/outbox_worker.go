package kafka

import (
	"log"
	"time"
	"workflow-service/internal/domain/repository"
)

type OutboxWorker struct {
	outboxRepo repository.OutboxRepository
	publisher  Producer
	topic      string
}

func NewOutboxWorker(o repository.OutboxRepository, p Producer, topic string) *OutboxWorker {
	return &OutboxWorker{outboxRepo: o, publisher: p, topic: topic}
}

func (w *OutboxWorker) Start() {
	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		events, err := w.outboxRepo.FindPending(10)
		if err != nil {
			log.Println("outbox fetch error : ", err)
			continue
		}

		for _, event := range events {
			if err := w.publisher.Publish(w.topic, event.Payload); err != nil {
				log.Println("kafka publish error : ", err)
				continue
			}
			_ = w.outboxRepo.MarkPublished(event.ID)
		}
	}
}
