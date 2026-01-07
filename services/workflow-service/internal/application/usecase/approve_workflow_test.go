package usecase

import (
	"testing"
	"workflow-service/internal/application/usecase/mocks"
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApproveWorkflow_SuccessOutboxEvent(t *testing.T) {
	// arrage
	repo := new(mocks.WorkflowRepository)
	cache := new(mocks.WorkflowCache)
	outbox := new(mocks.OutboxRepository)

	workflow := &entity.Workflow{
		ID:     1,
		Status: "CREATED",
	}

	repo.On("FindByID", "1").Return(workflow, nil)
	repo.On("Update", mock.Anything).Return(nil)

	cache.On("Delete", "1").Return(nil)
	cache.On("Set", mock.Anything).Return(nil)

	outbox.On("Create", mock.MatchedBy(func(e *entity.OutboxEvent) bool {
		return e.EventType == "WorkflowApproved" &&
			e.Status == "PENDING"
	})).Return(nil)

	uc := NewApproveWorkflowUsecase(repo, cache, outbox)

	// act
	err := uc.Approve("1")

	// assert
	assert.NoError(t, err)

	repo.AssertCalled(t, "Update", mock.Anything)
	outbox.AssertCalled(t, "Create", mock.Anything)
	cache.AssertCalled(t, "Delete", "1")
}
