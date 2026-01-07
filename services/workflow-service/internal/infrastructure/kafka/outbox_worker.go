package kafka

import (
	"log"
	"time"
	"workflow-service/internal/domain/repository"
)

type OutboxWorker struct {
	outboxRepo repository.OutboxRepository
	publisher  repository.EventPublisher
}

func NewOutboxWorker(o repository.OutboxRepository, p repository.EventPublisher) *OutboxWorker {
	return &OutboxWorker{outboxRepo: o, publisher: p}
}

func (w *OutboxWorker) Start() {
	for {
		events, err := w.outboxRepo.FindPending(10)
		if err != nil {
			log.Println("outbox fetch error : ", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, event := range events {
			err := w.publisher.Publish("workflow-events", event.Payload)
			if err != nil {
				log.Println("publish failed : ", err)
				continue
			}
			_ = w.outboxRepo.MarkPublished(event.ID)
		}
		time.Sleep(1 * time.Second)
	}
}
