package usecase

import (
	"log"
	"os"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
	"workflow-service/internal/infrastructure/kafka"
)

type CreateWorkflowUsecase struct {
	repo     repository.WorkflowRepository
	producer *kafka.Producer
	cache    repository.WorkflowCache
}

func NewCreateWorkflowUsecase(r repository.WorkflowRepository, p *kafka.Producer, c repository.WorkflowCache) *CreateWorkflowUsecase {
	return &CreateWorkflowUsecase{repo: r, producer: p, cache: c}
}

func (uc *CreateWorkflowUsecase) Create(name string) error {
	workflow := &entity.Workflow{
		Name:   name,
		Status: "CREATED",
	}

	// save to DB
	if err := uc.repo.Create(workflow); err != nil {
		log.Println("create error:", err)
		return err
	}

	// cache
	if err := uc.cache.Set(workflow); err != nil {
		log.Println("cache error : ", err)
		return err
	}

	event := entity.WorkflowEvent{
		WorkflowID:  workflow.ID,
		Name:        workflow.Name,
		Status:      workflow.Status,
		Type:        "WorkflowCreated",
		Description: "Workflow created by user",
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}

	if err := uc.producer.Publish(os.Getenv("KAFKA_TOPIC"), event); err != nil {
		log.Println("kafka error:", err)
		return err
	}

	return nil
}
