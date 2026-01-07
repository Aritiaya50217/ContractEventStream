package postgres

import (
	"workflow-service/internal/domain/entity"

	"gorm.io/gorm"
)

type OutboxRepoPG struct {
	db *gorm.DB
}

func NewOutboxRepoPG(db *gorm.DB) *OutboxRepoPG {
	return &OutboxRepoPG{db: db}
}

func (r *OutboxRepoPG) Create(event *entity.OutboxEvent) error {
	return r.db.Create(event).Error
}

func (r *OutboxRepoPG) FindPending(limit int) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent
	err := r.db.Where("status = ? ", "PENDING").Order("id ASC").Limit(limit).Find(&events).Error
	return events, err
}

func (r *OutboxRepoPG) MarkPublished(id uint) error {
	return r.db.Model(&entity.OutboxEvent{}).Where("id = ? ", id).Update("status", "PUBLISHED").Error
}
