package entity

import (
	"errors"
	"time"
)

type Workflow struct {
	ID        uint   `gorm:"primaryKey"` 
	Name      string `gorm:"type:varchar(255);not null"`
	Status    string `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *Workflow) Approve() error {
	if w.Status != "CREATED" {
		return errors.New("workflow cannot be approved")
	}
	w.Status = "APPROVED"
	w.UpdatedAt = time.Now()
	return nil
}
