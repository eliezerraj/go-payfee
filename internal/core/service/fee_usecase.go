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

// About add a script fee
func (s *WorkerService) AddScript(ctx context.Context, script *model.ScriptData) (*model.ScriptData, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("AddScript")
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("script: ",script).Msg("")

	//trace
	span := tracerProvider.Span(ctx, "service.AddScript")
	defer span.End()
	
	// prepare 
	key := "script:" + script.Script.Name
	value_json, err := json.Marshal(script.Script)
    if err != nil {
		return nil, err
    }

	// Put to the cache
	res, err := s.workerRepository.RedisClusterServer.Set(ctx, key, value_json)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, erro.ErrInsert
	}
	return script, nil
}

// About get a script fee
func (s *WorkerService) GetScript(ctx context.Context, script *model.ScriptData) (*model.Script, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetScript")
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("script: ",script).Msg("")

	// trace
	span := tracerProvider.Span(ctx, "service.GetScript")
	defer span.End()

	// Prepare
	key := "script:"+ script.Script.Name

	// Get to the cache
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

// About add a key
func (s *WorkerService) AddKey(ctx context.Context, fee *model.Fee) (*model.Fee, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("AddKey")
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("fee: ",fee).Msg("")

	// trace
	span := tracerProvider.Span(ctx, "service.AddKey")
	defer span.End()
	
	// prepare
	key := "fee:" + fee.Name
	valueReg := fee.Value

	// Put to the cache
	res, err := s.workerRepository.RedisClusterServer.Set(ctx, key, valueReg)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, erro.ErrInsert
	}
	return fee, nil
}

// About get a key
func (s *WorkerService) GetKey(ctx context.Context, fee *model.Fee) (*model.Fee, error){
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Msg("GetKey")
	childLogger.Info().Interface("trace-resquest-id", ctx.Value("trace-request-id")).Interface("fee: ",fee).Msg("")

	// Trace
	span := tracerProvider.Span(ctx, "service.GetKey")
	defer span.End()
	
	// Prepare
	key := "fee:" + fee.Name

	// Get to the cache
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