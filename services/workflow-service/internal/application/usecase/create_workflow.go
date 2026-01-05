package usecase

import (
	"encoding/json"
	"log"
	"os"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
)

type CreateWorkflowUsecase struct {
	repo repository.WorkflowRepository
	// producer *kafka.Producer
	publisher repository.EventPublisher
	cache     repository.WorkflowCache
}

func NewCreateWorkflowUsecase(r repository.WorkflowRepository, p repository.EventPublisher, c repository.WorkflowCache) *CreateWorkflowUsecase {
	return &CreateWorkflowUsecase{repo: r, publisher: p, cache: c}
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

	// if err := uc.producer.Publish(os.Getenv("KAFKA_TOPIC"), event); err != nil {
	// 	log.Println("kafka error:", err)
	// 	return err
	// }

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err := uc.publisher.Publish(os.Getenv("KAFKA_TOPIC"), payload); err != nil {
		return err
	}

	return nil
}
