package cache

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	//"github.com/go-payfee/internal/core"

	redis "github.com/redis/go-redis/v9"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var childLogger = log.With().Str("repository/cache", "Redis").Logger()

type CacheService struct {
	cache *redis.ClusterClient
}

func NewClusterCache(ctx context.Context, options *redis.ClusterOptions) *CacheService {
	childLogger.Debug().Msg("NewClusterCache")
	childLogger.Debug().Interface("option.Addrs: ", options.Addrs).Msg("")

	redisClient := redis.NewClusterClient(options)
	return &CacheService{
		cache: redisClient,
	}
}

func (s *CacheService) Ping(ctx context.Context) (string, error) {
	childLogger.Debug().Msg("Ping")

	status, err := s.cache.Ping(ctx).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *CacheService) Get(ctx context.Context, key string) (interface{}, error) {
	childLogger.Debug().Msg("Get")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Get")
	defer func() {
		root.Close(nil)
	}()

	res, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CacheService) Put(ctx context.Context, key string, value interface{}) error {
	childLogger.Debug().Msg("Put")
	//childLogger.Debug().Str("====> key : ",key).Interface("| Put : ",value).Msg("")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Put")
	defer func() {
		root.Close(nil)
	}()

	value_json, err := json.Marshal(value)
    if err != nil {
       return err
    }

	status := s.cache.Set(ctx, key, value_json, 0)

	return status.Err()
}

func (s *CacheService) SetCount(ctx context.Context, key string, valueReg string, value interface{}) (error) {
	childLogger.Debug().Msg("Count")

	_, root := xray.BeginSubsegment(ctx, "REDIS.HIncrByFloat")
	defer func() {
		root.Close(nil)
	}()

	_, err := s.cache.HIncrByFloat(ctx, key, valueReg, value.(float64)).Result()
	if err != nil {
		return err
	}

	//s.cache.PExpire(ctx, key, time.Minute * 1).Result()
	return nil
}

func (s *CacheService) GetCount(ctx context.Context, key string, valueReg string) (interface{}, error) {
	childLogger.Debug().Msg("GetCount")

	_, root := xray.BeginSubsegment(ctx, "REDIS.HGet")
	defer func() {
		root.Close(nil)
	}()

	res, err := s.cache.HGet(ctx, key, valueReg).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CacheService) AddKey(ctx context.Context, key string, valueReg interface{}) (bool, error) {
	childLogger.Debug().Msg("AddKey")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Set")
	defer func() {
		root.Close(nil)
	}()

	err := s.cache.Set(ctx, key, valueReg , 0).Err()
	if err != nil {
		childLogger.Error().Err(err).Msg("AddKey")
		return false, err
	}

	return true, nil
}

func (s *CacheService) GetKey(ctx context.Context,key string) (interface{}, error) {
	childLogger.Debug().Msg("GetKey")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Get")
	defer func() {
		root.Close(nil)
	}()

	result, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		childLogger.Error().Err(err).Msg("GetKey")
		return nil, err
	}

	return result, nil
}