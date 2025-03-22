package cache

import (
	go_core_cache "github.com/eliezerraj/go-core/cache/redis_cluster"

	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("component", "go-payfee").Str("package", "internal.adapter.cache").Logger()

type WorkerRepository struct {
	RedisClusterServer *go_core_cache.RedisClusterServer
}

// About create worker
func NewWorkerRepository(redisClusterServer *go_core_cache.RedisClusterServer) *WorkerRepository{
	childLogger.Info().Str("func","NewWorkerRepository").Send()

	return &WorkerRepository{
		RedisClusterServer: redisClusterServer,
	}
}
