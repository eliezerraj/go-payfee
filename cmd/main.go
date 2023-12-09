package main

import(
	"time"
	"context"
	"strings"
	"crypto/tls"
	"os"
	"strconv"
	"net"

	"github.com/go-payfee/internal/service"
	"github.com/go-payfee/internal/handler"
	"github.com/go-payfee/internal/repository/cache"
	"github.com/go-payfee/internal/core"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
    "github.com/aws/aws-sdk-go-v2/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	redis "github.com/redis/go-redis/v9"
)

var(
	logLevel 	= zerolog.DebugLevel
	version 	= "GO-PAYFEE 1.0"
	ctxTimeout  = 29 // Session TimeOut

	infoPod					core.InfoPod
	httpAppServerConfig 	core.HttpAppServer
	server					core.Server

	envCacheCluster			redis.ClusterOptions

	noAZ		=	true // set only if you get to split the xray trace per AZ
)

func getEnv() {
	log.Debug().Msg("getEnv")

	if os.Getenv("API_VERSION") !=  "" {
		infoPod.ApiVersion = os.Getenv("API_VERSION")
	}
	if os.Getenv("POD_NAME") !=  "" {
		infoPod.PodName = os.Getenv("POD_NAME")
	}

	if os.Getenv("PORT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("PORT"))
		server.Port = intVar
	}

	if os.Getenv("REDIS_CLUSTER_ADDRESS") !=  "" {
		infoPod.RedisAddr = os.Getenv("REDIS_CLUSTER_ADDRESS")
		envCacheCluster.Addrs = strings.Split(os.Getenv("REDIS_CLUSTER_ADDRESS"), ",") 
	}

	if os.Getenv("NO_AZ") == "false" {	
		noAZ = false
	} else {
		noAZ = true
	}
}

func init(){
	log.Debug().Msg("init")

	server.Port = 5004
	server.ReadTimeout = 60
	server.WriteTimeout = 60
	server.IdleTimeout = 60
	server.CtxTimeout = 60

	envCacheCluster.Username = ""
	envCacheCluster.Password = ""
	envCacheCluster.Addrs = strings.Split("clustercfg.memdb-arch.vovqz2.memorydb.us-east-2.amazonaws.com:6379", ",")

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error().Err(err).Msg("Error to get the POD IP address !!!")
		os.Exit(3)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				infoPod.IPAddress = ipnet.IP.String()
			}
		}
	}
	infoPod.OSPID = strconv.Itoa(os.Getpid())

	getEnv()

	// Get AZ only if localtest is true
	if (noAZ != true) {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Error().Err(err).Msg("ERRO FATAL get Context !!!")
			os.Exit(3)
		}
		client := imds.NewFromConfig(cfg)
		response, err := client.GetInstanceIdentityDocument(context.TODO(), &imds.GetInstanceIdentityDocumentInput{})
		if err != nil {
			log.Error().Err(err).Msg("Unable to retrieve the region from the EC2 instance !!!")
			os.Exit(3)
		}
		infoPod.AvailabilityZone = response.AvailabilityZone	
	} else {
		infoPod.AvailabilityZone = "LOCALHOST_NO_AZ"
	}
}

func main() {
	log.Debug().Msg("main")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration( ctxTimeout ) * time.Second)
	defer cancel()

	if !strings.Contains(envCacheCluster.Addrs[0], "127.0.0.1") {
		envCacheCluster.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	cacheRedis := cache.NewClusterCache(ctx, &envCacheCluster)
	_, err := cacheRedis.Ping(ctx)
	if err != nil{
		log.Error().Err(err).Msg("Erro na abertura do Redis")
		panic(err)
	}

	log.Debug().Msg("Redis Ping Sucessful !!!")

	workerService := service.NewRBACService(cacheRedis)
	httpWorkerAdapter := handler.NewHttpWorkerAdapter(workerService)

	httpAppServerConfig.InfoPod = &infoPod
	httpAppServerConfig.Server = server
	httpServer := handler.NewHttpAppServer(httpAppServerConfig)

	httpServer.StartHttpAppServer(ctx, httpWorkerAdapter)
}