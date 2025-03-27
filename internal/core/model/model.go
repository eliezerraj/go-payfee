package model

import (
	redis "github.com/redis/go-redis/v9"
	go_core_observ "github.com/eliezerraj/go-core/observability" 
)

type AppServer struct {
	InfoPod 		*InfoPod 					`json:"info_pod"`
	Server     		*Server     				`json:"server"`
	RedisAddress	*string						`json:"redis_address"`
	ConfigOTEL		*go_core_observ.ConfigOTEL	`json:"otel_config"`
	DatabaseRedis	*DatabaseRedis  			`json:"redis_cluster"`	
}

type InfoPod struct {
	PodName				string 	`json:"pod_name"`
	ApiVersion			string 	`json:"version"`
	OSPID				string 	`json:"os_pid"`
	IPAddress			string 	`json:"ip_address"`
	AvailabilityZone 	string 	`json:"availabilityZone"`
	IsAZ				bool   	`json:"is_az"`
	Env					string `json:"enviroment,omitempty"`
	AccountID			string `json:"account_id,omitempty"`
}

type Server struct {
	Port 			int `json:"port"`
	ReadTimeout		int `json:"readTimeout"`
	WriteTimeout	int `json:"writeTimeout"`
	IdleTimeout		int `json:"idleTimeout"`
	CtxTimeout		int `json:"ctxTimeout"`
}

type MessageRouter struct {
	Message			string `json:"message"`
}

type DatabaseRedis struct {
	RedisAddress	string		`json:"redis_address"`
	RedisOptions	redis.ClusterOptions
}

type ScriptData struct {
    Script		Script 	`redis:"script" json:"script"`
}

type Script struct {
    Name 		string  `redis:"name" json:"name"`
    Description string   `redis:"description" json:"description"`
	Fee		    []string `redis:"fee" json:"fee"`
}

type Fee struct {
    Name 		string  `redis:"name" json:"name"`
	Value		float64  `redis:"value" json:"value"`
}