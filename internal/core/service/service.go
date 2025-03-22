package service

import(
	"github.com/go-payfee/internal/adapter/cache"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("component","go-payfee").Str("package","internal.core.service").Logger()

type WorkerService struct {
	workerRepository *cache.WorkerRepository
}

func NewWorkerService(workerRepository *cache.WorkerRepository) *WorkerService{
	childLogger.Info().Str("func","NewWorkerService").Send()

	return &WorkerService{
		workerRepository: workerRepository,
	}
}