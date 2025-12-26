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
	cache    repository.WorkflowCache
}

func NewApproveWorkflowUsecase(r repository.WorkflowRepository, p *kafka.Producer, c repository.WorkflowCache) *ApproveWorkflowUsecase {
	return &ApproveWorkflowUsecase{repo: r, producer: p, cache: c}
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

	// refresh cache (delete + set)
	if err := uc.cache.Delete(id); err != nil {
		log.Println("cache delete error:", err)
		return err
	}

	if err := uc.cache.Set(workflow); err != nil {
		log.Println("cache set error:", err)
		return err
	}

	event := entity.WorkflowEvent{
		WorkflowID:  workflow.ID,
		Name:        workflow.Name,
		Status:      workflow.Status,
		Type:        "WorkflowApproved",
		Description: "Workflow approved by user",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.producer.Publish(os.Getenv("KAFKA_TOPIC"), event); err != nil {
		log.Println("kafka error : ", err)
		return err
	}

	return nil
}
