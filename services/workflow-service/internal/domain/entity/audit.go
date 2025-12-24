package entity

import "time"

// AuditLog เก็บ History
type AuditLog struct {
	ID          uint `gorm:"primaryKey"`
	WorkflowID  uint
	EventType   string
	Status      string
	UserID      string
	UpdatedAt   time.Time
	Description string
}
