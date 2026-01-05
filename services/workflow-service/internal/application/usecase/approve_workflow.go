package usecase

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
)

type ApproveWorkflowUsecase struct {
	repo repository.WorkflowRepository
	// producer *kafka.Producer
	publisher repository.EventPublisher
	cache     repository.WorkflowCache
}

func NewApproveWorkflowUsecase(r repository.WorkflowRepository, p repository.EventPublisher, c repository.WorkflowCache) *ApproveWorkflowUsecase {
	return &ApproveWorkflowUsecase{repo: r, publisher: p, cache: c}
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
	_ = uc.cache.Delete(id)
	_ = uc.cache.Set(workflow)

	event := entity.WorkflowEvent{
		WorkflowID:  workflow.ID,
		Name:        workflow.Name,
		Status:      workflow.Status,
		Type:        "WorkflowApproved",
		Description: "Workflow approved by user",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// if err := uc.producer.Publish(os.Getenv("KAFKA_TOPIC"), event); err != nil {
	// 	log.Println("kafka error : ", err)
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
