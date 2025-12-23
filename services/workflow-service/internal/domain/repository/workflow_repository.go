package repository

import "workflow-service/internal/domain/entity"

type WorkflowRepository interface {
	Create(workflow *entity.Workflow) error
	FindByID(id string) (*entity.Workflow, error)
	Update(workflow *entity.Workflow) error
}
