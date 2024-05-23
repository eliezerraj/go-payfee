package main

import(
	"time"
	"context"
	
	"github.com/go-payfee/internal/service"
	"github.com/go-payfee/internal/handler"
	"github.com/go-payfee/internal/util"
	"github.com/go-payfee/internal/repository/cache"
	"github.com/go-payfee/internal/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var(
	logLevel 	= zerolog.DebugLevel
	appServer	core.AppServer
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	infoPod, server := util.GetInfoPod()
	database := util.GetDatabaseEnv()
	configOTEL := util.GetOtelEnv()

	appServer.InfoPod = &infoPod
	appServer.Database = &database
	appServer.Server = &server
	appServer.ConfigOTEL = &configOTEL
	appServer.RedisAddress = &database.RedisAddress
}

func main() {
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Msg("main")
	log.Debug().Msg("----------------------------------------------------")
	log.Debug().Interface("appServer.InfoPod :",appServer.InfoPod).Msg("")
	log.Debug().Interface("appServer.RedisAddress :",appServer.RedisAddress).Msg("")
	log.Debug().Interface("appServer.Server :",appServer.Server).Msg("")
	log.Debug().Interface("appServer.ConfigOTEL :",appServer.ConfigOTEL).Msg("")
	log.Debug().Msg("----------------------------------------------------")

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( appServer.Server.ReadTimeout ) * time.Second)
	defer cancel()

	cacheRedis := cache.NewClusterCache(ctx, 
										&appServer.Database.RedisOptions)
	_, err := cacheRedis.Ping(ctx)
	if err != nil{
		log.Error().Err(err).Msg("Erro na abertura do Redis")
		panic(err)
	}
	log.Debug().Msg(" ===> Redis Ping Sucessful !!! <===")

	workerService := service.NewRBACService(cacheRedis)

	httpWorkerAdapter := handler.NewHttpWorkerAdapter(workerService)
	httpServer := handler.NewHttpAppServer(appServer.Server)
	
	httpServer.StartHttpAppServer(	ctx, 
									&httpWorkerAdapter,
									&appServer)
}