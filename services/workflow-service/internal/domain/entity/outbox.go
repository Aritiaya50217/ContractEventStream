package entity

import "time"

type OutboxEvent struct {
	ID        uint   `gorm:"primaryKey"`
	Aggregate string `gorm:"index"` // workflow
	EventType string // WorkflowApproved
	Payload   []byte `gorm:"type:jsonb"`
	Status    string `gorm:"index"` // PENDING, PUBLISHED
	CreatedAt time.Time
}
