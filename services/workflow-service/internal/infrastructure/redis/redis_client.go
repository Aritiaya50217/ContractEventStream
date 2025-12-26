package redis

import (
	"context"
	"log"
	"workflow-service/internal/config"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisClient() *redis.Client {
	cfg := config.LoadRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   cfg.DB,
	})

	// ping test
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	log.Println("Redis connected:", cfg.Addr)
	return client
}
