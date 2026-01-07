package mocks

import (
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type OutboxRepository struct {
	mock.Mock
}

func (m *OutboxRepository) Create(e *entity.OutboxEvent) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *OutboxRepository) FindPending(limit int) ([]entity.OutboxEvent, error) {
	args := m.Called(limit)
	return args.Get(0).([]entity.OutboxEvent), args.Error(1)

}

func (m *OutboxRepository) MarkPublished(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
