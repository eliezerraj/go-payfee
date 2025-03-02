package service

import(
	//"fmt"
	"context"
	"errors"
	"strconv"
	"encoding/json"

	"github.com/go-payfee/internal/core/model"
	"github.com/go-payfee/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
)

var tracerProvider go_core_observ.TracerProvider

func (s *WorkerService) AddScript(ctx context.Context, script *model.ScriptData) (*model.ScriptData, error){
	childLogger.Debug().Msg("AddScript")
	childLogger.Debug().Interface("script: ",script).Msg("")

	span := tracerProvider.Span(ctx, "service.AddScript")
	defer span.End()
	
	// Put to the cache
	key := "script:" + script.Script.Name

	value_json, err := json.Marshal(script.Script)
    if err != nil {
		return nil, err
    }

	res, err := s.workerRepository.RedisClusterServer.Set(ctx, key, value_json)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, erro.ErrInsert
	}
	return script, nil
}

func (s *WorkerService) GetScript(ctx context.Context, script *model.ScriptData) (*model.Script, error){
	childLogger.Debug().Msg("GetScript")
	childLogger.Debug().Interface("script: ",script).Msg("")

	span := tracerProvider.Span(ctx, "service.GetScript")
	defer span.End()

	// Get to the cache
	key := "script:"+ script.Script.Name

	res, err := s.workerRepository.RedisClusterServer.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	
	var script_assert model.Script
	err = json.Unmarshal([]byte(res.(string)), &script_assert)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &script_assert, nil
}

func (s *WorkerService) AddKey(ctx context.Context, fee *model.Fee) (*model.Fee, error){
	childLogger.Debug().Msg("AddKey")
	childLogger.Debug().Interface("fee: ",fee).Msg("")

	span := tracerProvider.Span(ctx, "service.AddKey")
	defer span.End()
	
	// Put to the cache
	key := "fee:" + fee.Name
	valueReg := fee.Value

	res, err := s.workerRepository.RedisClusterServer.Set(ctx, key, valueReg)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, erro.ErrInsert
	}
	return fee, nil
}

func (s *WorkerService) GetKey(ctx context.Context, fee *model.Fee) (*model.Fee, error){
	childLogger.Debug().Msg("GetKey")
	childLogger.Debug().Interface("fee: ",fee).Msg("")

	span := tracerProvider.Span(ctx, "service.GetKey")
	defer span.End()
	
	// Get to the cache
	key := "fee:" + fee.Name
	
	res, err := s.workerRepository.RedisClusterServer.Get(ctx, key)
	if err != nil {
		return nil, erro.ErrNotFound
	}

	value_assert, err := strconv.ParseFloat(res.(string), 64)
	if err != nil {
		return nil, err
	}

	res_fee := model.Fee{}
	res_fee.Name = fee.Name
	res_fee.Value = value_assert

	return &res_fee, nil
}