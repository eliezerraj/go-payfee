package configuration

import(
	"os"
	"strings"

	"crypto/tls"
	"github.com/joho/godotenv"
	"github.com/go-payfee/internal/core/model"
	
	redis "github.com/redis/go-redis/v9"
)

// About get redis env var
func GetRedisEnv() model.DatabaseRedis {
	childLogger.Info().Str("func","GetRedisEnv").Send()

	err := godotenv.Load(".env")
	if err != nil {
		childLogger.Info().Err(err).Send()
	}
	
	var databaseRedis		model.DatabaseRedis
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