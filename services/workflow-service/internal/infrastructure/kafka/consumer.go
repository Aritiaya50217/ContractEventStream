package kafka

import (
	"context"
	"encoding/json"
	"log"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/infrastructure/postgres"

	"github.com/segmentio/kafka-go"
)

type WorkflowConsumer struct {
	reader    *kafka.Reader
	auditRepo *postgres.AuditRepoPG
}

func NewWorkflowConsumer(r *kafka.Reader, a *postgres.AuditRepoPG) *WorkflowConsumer {
	return &WorkflowConsumer{reader: r, auditRepo: a}
}

func (c *WorkflowConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("kafka read error : ", err)
			continue
		}
		log.Println("Kafka message : ", string(msg.Value))

		var event entity.WorkflowEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		c.saveAudit(event)
		c.handleEvent(event)
	}
}

func (c *WorkflowConsumer) saveAudit(event entity.WorkflowEvent) {
	audit := entity.AuditLog{
		WorkflowID:  event.WorkflowID,
		EventType:   event.Type,
		Status:      event.Status,
		UserID:      event.UserID,
		UpdatedAt:   event.UpdatedAt,
		Description: event.Description,
	}

	if err := c.auditRepo.Create(&audit); err != nil {
		log.Println("Save audit error : ", err)
	}
}

func (c *WorkflowConsumer) handleEvent(event entity.WorkflowEvent) {
	switch event.Type {
	case "WorkflowCreated":
		log.Printf("Workflow Created: %v Status: %v\n", event.WorkflowID, event.Status)
	case "WorkflowApproved":
		log.Printf("Workflow Approved: %v Status: %v\n", event.WorkflowID, event.Status)
	default:
		log.Printf("Unknown event type: %s\n", event.Type)
	}
}
