package usecase

import (
	"encoding/json"
	"log"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
)

type CreateWorkflowUsecase struct {
	repo       repository.WorkflowRepository
	cache      repository.WorkflowCache
	outboxRepo repository.OutboxRepository
}

func NewCreateWorkflowUsecase(r repository.WorkflowRepository, c repository.WorkflowCache, o repository.OutboxRepository) *CreateWorkflowUsecase {
	return &CreateWorkflowUsecase{repo: r, cache: c, outboxRepo: o}
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

	// create outbox event (แทน kafka)
	event := entity.WorkflowEvent{
		WorkflowID:  workflow.ID,
		Name:        workflow.Name,
		Status:      workflow.Status,
		Type:        "WorkflowCreated",
		Description: "Workflow created by user",
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outbox := &entity.OutboxEvent{
		Aggregate: "workflow",
		EventType: "WorkflowCreated",
		Payload:   payload,
		Status:    "PENDING",
	}

	if err := uc.outboxRepo.Create(outbox); err != nil {
		return err
	}

	// cache
	if err := uc.cache.Set(workflow); err != nil {
		log.Println("cache error : ", err)
		return err
	}

	return nil
}
