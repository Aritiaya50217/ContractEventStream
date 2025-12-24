package entity

import "time"

type WorkflowEvent struct {
	WorkflowID  uint      `json:"workflow_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
