package mocks

import (
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type WorkflowCache struct {
	mock.Mock
}

func (m *WorkflowCache) Get(id string) (*entity.Workflow, error) {
	args := m.Called(id)
	if w := args.Get(0); w != nil {
		return w.(*entity.Workflow), args.Error(1)
	}
	return nil, args.Error(1)

}

func (m *WorkflowCache) Set(workflow *entity.Workflow) error {
	args := m.Called(workflow)
	return args.Error(0)
}

func (m *WorkflowCache) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
