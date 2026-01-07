package repository

import "workflow-service/internal/domain/entity"

type OutboxRepository interface {
	Create(event *entity.OutboxEvent) error
	FindPending(limit int) ([]entity.OutboxEvent, error)
	MarkPublished(id uint) error
}
