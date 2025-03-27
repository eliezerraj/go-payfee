package main

import(
	"time"
	"context"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-payfee/internal/infra/configuration"
	"github.com/go-payfee/internal/core/model"
	"github.com/go-payfee/internal/core/service"
	"github.com/go-payfee/internal/infra/server"
	"github.com/go-payfee/internal/adapter/api"
	"github.com/go-payfee/internal/adapter/cache"
	go_core_cache "github.com/eliezerraj/go-core/cache/redis_cluster"
)

var(
	logLevel 			= 	zerolog.InfoLevel // zerolog.InfoLevel zerolog.DebugLevel
	appServer			model.AppServer
	redisClusterServer 	go_core_cache.RedisClusterServer
	childLogger = log.With().Str("component","go-payfee").Str("package", "main").Logger()
)

// About initialize the enviroment var
func init(){
	childLogger.Info().Str("func","init").Send()

	zerolog.SetGlobalLevel(logLevel)

	infoPod, server := configuration.GetInfoPod()
	configOTEL 		:= configuration.GetOtelEnv()
	databaseRedis 	:= configuration.GetRedisEnv() 

	appServer.InfoPod = &infoPod
	appServer.Server = &server
	appServer.ConfigOTEL = &configOTEL
	appServer.DatabaseRedis = &databaseRedis
}


// About main
func main (){
	childLogger.Info().Str("func","main").Interface("appServer",appServer).Send()

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()

	// Open Database
	cacheRedis := redisClusterServer.NewClusterCache(&appServer.DatabaseRedis.RedisOptions)
	_, err := cacheRedis.Ping(ctx)
	if err != nil{
		childLogger.Error().Err(err).Msg("erro open redis")
		panic(err)
	}
	childLogger.Info().Msg("ping redis successful !!!")

	// wire	
	workerRepository := cache.NewWorkerRepository(cacheRedis)
	workerService := service.NewWorkerService(workerRepository)
	httpRouters := api.NewHttpRouters(workerService)
	httpServer := server.NewHttpAppServer(appServer.Server)

	// start server
	httpServer.StartHttpAppServer(ctx, &httpRouters, &appServer)
}