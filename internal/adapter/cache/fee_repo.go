package cache

import (
	go_core_cache "github.com/eliezerraj/go-core/cache/redis_cluster"

	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("adapter", "cache").Logger()

type WorkerRepository struct {
	RedisClusterServer *go_core_cache.RedisClusterServer
}

func NewWorkerRepository(redisClusterServer *go_core_cache.RedisClusterServer) *WorkerRepository{
	childLogger.Debug().Msg("NewWorkerRepository")

	return &WorkerRepository{
		RedisClusterServer: redisClusterServer,
	}
}
