package usecase

import (
	"os"
	"testing"
	"workflow-service/internal/application/usecase/mocks"
	"workflow-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApproveWorkflow_PublishKafkaEvent(t *testing.T) {
	// arrange
	os.Setenv("KAFKA_TOPIC", "workflow-events")
	defer os.Unsetenv("KAFKA_TOPIC")

	repo := new(mocks.WorkflowRepository)
	cache := new(mocks.WorkflowCache)
	publisher := new(mocks.EventPublisher)

	workflow := &entity.Workflow{
		ID:     1,
		Status: "CREATED",
	}

	repo.On("FindByID", "1").Return(workflow, nil)
	repo.On("Update", mock.Anything).Return(nil)
	cache.On("Delete", "1").Return(nil)
	cache.On("Set", mock.Anything).Return(nil)

	publisher.On("Publish", "workflow-events", mock.Anything).Return(nil)

	uc := NewApproveWorkflowUsecase(repo, publisher, cache)

	// act
	err := uc.Approve("1")

	// assert
	assert.NoError(t, err)
	publisher.AssertCalled(t, "Publish", "workflow-events", mock.Anything)
}
