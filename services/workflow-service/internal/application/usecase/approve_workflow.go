package usecase

import "workflow-service/internal/domain/repository"

type ApproveWorkflowUsecase struct {
	repo repository.WorkflowRepository
}

func NewApproveWorkflowUsecase(r repository.WorkflowRepository) *ApproveWorkflowUsecase {
	return &ApproveWorkflowUsecase{repo: r}
}

func (uc *ApproveWorkflowUsecase) Approve(id string) error {
	workflow, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	if err := workflow.Approve(); err != nil {
		return err
	}

	return uc.repo.Update(workflow)
}
