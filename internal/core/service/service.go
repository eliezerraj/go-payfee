package service

import(
	"github.com/go-payfee/internal/adapter/cache"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("core", "service").Logger()

type WorkerService struct {
	workerRepository *cache.WorkerRepository
}

func NewWorkerService(workerRepository *cache.WorkerRepository) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepository: workerRepository,
	}
}