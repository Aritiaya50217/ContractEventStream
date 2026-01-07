package usecase

import (
	"testing"
	"workflow-service/internal/application/usecase/mocks"
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWorkflow_SaveOutboxEvent(t *testing.T) {
	// arrange
	repo := new(mocks.WorkflowRepository)
	cache := new(mocks.WorkflowCache)
	outbox := new(mocks.OutboxRepository)

	repo.On("Create", mock.AnythingOfType("*entity.Workflow")).Return(nil)

	cache.On("Set", mock.AnythingOfType("*entity.Workflow")).Return(nil)

	outbox.On("Create", mock.MatchedBy(func(e *entity.OutboxEvent) bool {
		return e.EventType == "WorkflowCreated" && e.Status == "PENDING"
	})).Return(nil)

	outbox.On("FindPending", mock.Anything).Return([]entity.OutboxEvent{}, nil)

	outbox.On("MarkPublished", mock.Anything).Return(nil)

	uc := NewCreateWorkflowUsecase(repo, cache, outbox)

	// Act
	err := uc.Create("Test Workflow")

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "Create", mock.Anything)
	outbox.AssertCalled(t, "Create", mock.Anything)
}
