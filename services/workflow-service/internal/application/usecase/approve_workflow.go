package usecase

import (
	"encoding/json"
	"log"
	"time"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
)

type ApproveWorkflowUsecase struct {
	repo repository.WorkflowRepository
	// producer *kafka.Producer
	// publisher repository.EventPublisher
	cache      repository.WorkflowCache
	outboxRepo repository.OutboxRepository
}

func NewApproveWorkflowUsecase(r repository.WorkflowRepository, c repository.WorkflowCache, o repository.OutboxRepository) *ApproveWorkflowUsecase {
	return &ApproveWorkflowUsecase{repo: r, cache: c, outboxRepo: o}
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

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := &entity.OutboxEvent{
		Aggregate: "workflow",
		EventType: "WorkflowApproved",
		Payload:   payload,
		Status:    "PENDING",
	}

	if err := uc.outboxRepo.Create(outboxEvent); err != nil {
		return err
	}
	return nil
}
