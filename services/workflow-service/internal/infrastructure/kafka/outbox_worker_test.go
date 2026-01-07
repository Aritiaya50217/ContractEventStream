package kafka

import (
	"testing"
	"time"
	"workflow-service/internal/application/usecase/mocks"
	"workflow-service/internal/domain/entity"
)

func TestOutboxWorker_PublishesEvent(t *testing.T) {
	// mocks
	outboxRepo := new(mocks.OutboxRepository)
	publisher := new(mocks.EventPublisher)

	event := entity.OutboxEvent{
		ID:        1,
		EventType: "WorkflowCreated",
		Status:    "PENDING",
		Payload:   []byte(`{"workflow_id":1}`),
	}

	// กำหนด behavior ของ mock
	outboxRepo.On("FindPending", 10).Return([]entity.OutboxEvent{event}, nil).Once()
	outboxRepo.On("MarkPublished", uint(1)).Return(nil).Once()
	publisher.On("Publish", "workflow-events", event.Payload).Return(nil).Once()

	// สร้าง worker
	worker := NewOutboxWorker(outboxRepo, publisher)

	// run worker ใน goroutine
	go worker.Start()
	time.Sleep(100 * time.Millisecond) // ให้ worker loop 1 ครั้ง

	// assert
	publisher.AssertCalled(t, "Publish", "workflow-events", event.Payload)
	outboxRepo.AssertCalled(t, "MarkPublished", uint(1))
}
