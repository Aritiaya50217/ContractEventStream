package usecase

import (
	"log"
	"os"
	"time"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
	"workflow-service/internal/infrastructure/kafka"
)

type ApproveWorkflowUsecase struct {
	repo     repository.WorkflowRepository
	producer *kafka.Producer
}

func NewApproveWorkflowUsecase(r repository.WorkflowRepository, p *kafka.Producer) *ApproveWorkflowUsecase {
	return &ApproveWorkflowUsecase{repo: r, producer: p}
}

func (uc *ApproveWorkflowUsecase) Approve(id string) error {
	workflow, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err := workflow.Approve(); err != nil {
		return err
	}

	// save to DB
	if err := uc.repo.Update(workflow); err != nil {
		log.Println("update error : ", err)
		return err
	}

	event := entity.WorkflowEvent{
		ID:        workflow.ID,
		Name:      workflow.Name,
		Status:    workflow.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.producer.Publish(os.Getenv("KAFKA_APPROVE_TOPIC"), event); err != nil {
		log.Println("kafka error : ", err)
		return err
	}

	return nil
}
