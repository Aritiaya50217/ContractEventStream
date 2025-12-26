package usecase

import (
	"log"
	"workflow-service/internal/domain/entity"
	"workflow-service/internal/domain/repository"
)

type GetWorkflowUsecase struct {
	repo  repository.WorkflowRepository
	cache repository.WorkflowCache
}

func NewGetWorkflowUsecase(repo repository.WorkflowRepository, cache repository.WorkflowCache) *GetWorkflowUsecase {
	return &GetWorkflowUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *GetWorkflowUsecase) GetByID(id string) (*entity.Workflow, error) {
	// cache
	if workflow, err := uc.cache.Get(id); err != nil {
		return workflow, err
	}

	// db
	workflow, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// set cache
	if err := uc.cache.Set(workflow); err != nil {
		log.Println("cache error : ", err)
		return nil, err
	}

	return workflow, nil
}
