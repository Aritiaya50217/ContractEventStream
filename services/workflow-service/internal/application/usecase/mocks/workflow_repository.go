package mocks

import (
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type WorkflowRepository struct {
	mock.Mock
}

func (m *WorkflowRepository) FindByID(id string) (*entity.Workflow, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Workflow), args.Error(1)
}

func (m *WorkflowRepository) Update(w *entity.Workflow) error {
	args := m.Called(w)
	return args.Error(0)
}

func (m *WorkflowRepository) Create(w *entity.Workflow) error {
	args := m.Called(w)
	return args.Error(0)
}
