package util

import(
	"os"
	"strings"

	"crypto/tls"
	"github.com/joho/godotenv"
	"github.com/go-payfee/internal/core"
	redis "github.com/redis/go-redis/v9"
)

func GetDatabaseEnv() core.DatabaseRedis {
	childLogger.Debug().Msg("GetDatabaseEnv")

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Msg("No .env File !!!!")
	}
	
	var databaseRedis		core.DatabaseRedis
	var envCacheCluster		redis.ClusterOptions

	envCacheCluster.Username = ""
	envCacheCluster.Password = ""

	if os.Getenv("REDIS_CLUSTER_ADDRESS") !=  "" {
		databaseRedis.RedisAddress = os.Getenv("REDIS_CLUSTER_ADDRESS")
		envCacheCluster.Addrs = strings.Split(os.Getenv("REDIS_CLUSTER_ADDRESS"), ",") 
	}

	// Just for local test
	if !strings.Contains(envCacheCluster.Addrs[0], "127.0.0.1") {
		envCacheCluster.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	databaseRedis.RedisOptions = envCacheCluster

	return databaseRedis
}