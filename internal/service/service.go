package service

import (
		"context"
		"errors"
		"strconv"
		"encoding/json"
		"fmt"

		"github.com/mitchellh/mapstructure"
		"github.com/rs/zerolog/log"
		"github.com/go-payfee/internal/repository/cache"
		"github.com/go-payfee/internal/core"
		"github.com/go-payfee/internal/lib"
)

var childLogger = log.With().Str("service", "RedisService").Logger()

type RedisService struct {
	cacheRedis	*cache.CacheService
}

func NewRBACService(cacheRedis *cache.CacheService) *RedisService{
	childLogger.Debug().Msg("NewRBACService")

	return &RedisService{
		cacheRedis: cacheRedis,
	}
}

func (s *RedisService) AddScript(ctx context.Context, script core.ScriptData) (bool, error) {
	childLogger.Debug().Msg("addScript")

	span := lib.Span(ctx, "service.addScript")
	defer span.End()

	// Put to the cache
	key := "script:" + script.Script.Name

	err := s.cacheRedis.Put(ctx, key, script.Script)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return false, err
	}

	return true, nil
}

func (s *RedisService) GetScript(ctx context.Context, script core.ScriptData) (*core.Script, error) {
	childLogger.Debug().Msg("GetScript")

	span := lib.Span(ctx, "service.GetScript")
	defer span.End()

	// Get to the cache
	key := "script:"+ script.Script.Name

	res, err := s.cacheRedis.Get(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return nil, err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(fmt.Sprint(res)), &jsonMap)

	var script_assert core.Script
	err = mapstructure.Decode(jsonMap, &script_assert)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	return &script_assert, nil
}

func (s *RedisService) AddKey(ctx context.Context, fee core.Fee) (bool, error) {
	childLogger.Debug().Msg("AddKey")

	span := lib.Span(ctx, "service.AddKey")
	defer span.End()

	// Put to the cache
	key := "fee:" + fee.Name
	valueReg := fee.Value

	res, err := s.cacheRedis.AddKey(ctx, key, valueReg)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return false, err
	}

	return res, nil
}

func (s *RedisService) GetKey(ctx context.Context, fee core.Fee) (*core.Fee, error) {
	childLogger.Debug().Msg("GetKey")

	span := lib.Span(ctx, "service.GetKey")
	defer span.End()

	// Get to the cache
	key := "fee:" + fee.Name
	
	res, err := s.cacheRedis.GetKey(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return nil, err
	}

	value_assert, err := strconv.ParseFloat(res.(string), 64)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return nil, err
	}

	res_fee := core.Fee{}
	res_fee.Name = fee.Name
	res_fee.Value = value_assert

	return &res_fee, nil
}