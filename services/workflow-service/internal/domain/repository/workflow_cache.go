package repository

import (
	"workflow-service/internal/domain/entity"
)

type WorkflowCache interface {
	Get(id string) (*entity.Workflow, error)
	Set(wf *entity.Workflow) error
	Delete(id string) error
}
