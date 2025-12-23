package postgres

import (
	"workflow-service/internal/domain/entity"

	"gorm.io/gorm"
)

type WorkflowRepoPg struct {
	db *gorm.DB
}

func NewWorkflowRepoPg(db *gorm.DB) *WorkflowRepoPg {
	return &WorkflowRepoPg{db: db}
}

func (r *WorkflowRepoPg) Create(workflow *entity.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowRepoPg) FindByID(id string) (*entity.Workflow, error) {
	var wf entity.Workflow
	if err := r.db.First(&wf, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *WorkflowRepoPg) Update(workflow *entity.Workflow) error {
	return r.db.Save(workflow).Error
}
