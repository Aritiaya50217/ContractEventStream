package config

import (
	"fmt"
	"os"
)

type RedisConfig struct {
	Addr string
	DB   int
	TTL  int
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		DB:  0,
		TTL: 60,
	}
}
