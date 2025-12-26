package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"workflow-service/internal/domain/entity"

	"github.com/go-redis/redis/v8"
)

type WorkflowCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewWorkflowCache(client *redis.Client) *WorkflowCache {
	return &WorkflowCache{
		client: client,
		ttl:    5 * time.Minute,
	}
}

func (c *WorkflowCache) key(id string) string {
	return fmt.Sprintf("workflow:%v", id)
}

func (c *WorkflowCache) Get(id string) (*entity.Workflow, error) {
	ctx := context.Background()

	val, err := c.client.Get(ctx, c.key(id)).Result()
	if err != nil {
		return nil, err
	}

	var workflow entity.Workflow
	if err := json.Unmarshal([]byte(val), &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (c *WorkflowCache) Set(workflow *entity.Workflow) error {
	ctx := context.Background()

	value, err := json.Marshal(workflow)
	if err != nil {
		return err
	}
	id := strconv.FormatUint(uint64(workflow.ID), 10)
	return c.client.Set(ctx, c.key(id), value, c.ttl).Err()
}

func (c *WorkflowCache) Delete(id string) error {
	ctx := context.Background()
	return c.client.Del(ctx, c.key(id)).Err()
}
