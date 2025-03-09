package server

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/go-payfee/internal/core/model"
	"github.com/go-payfee/internal/adapter/api"

	go_core_observ "github.com/eliezerraj/go-core/observability"  

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/eliezerraj/go-core/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var childLogger = log.With().Str("handler", "api").Logger()
var core_middleware middleware.ToolsMiddleware
var tracerProvider go_core_observ.TracerProvider
var infoTrace go_core_observ.InfoTrace

type HttpServer struct {
	httpServer	*model.Server
}

func NewHttpAppServer(httpServer *model.Server) HttpServer {
	return HttpServer{httpServer: httpServer }
}

// About start the http server
func (h HttpServer) StartHttpAppServer(	ctx context.Context, 
										httpRouters *api.HttpRouters,
										appServer *model.AppServer) {
	childLogger.Info().Msg("StartHttpAppServer")
			
	// otel
	childLogger.Info().Str("OTEL_EXPORTER_OTLP_ENDPOINT :", appServer.ConfigOTEL.OtelExportEndpoint).Msg("")
	
	infoTrace.PodName = appServer.InfoPod.PodName
	infoTrace.PodVersion = appServer.InfoPod.ApiVersion
	infoTrace.ServiceType = "k8-workload"
	infoTrace.Env = appServer.InfoPod.Env
	infoTrace.AccountID = appServer.InfoPod.AccountID

	tp := tracerProvider.NewTracerProvider(	ctx, 
											appServer.ConfigOTEL, 
											&infoTrace)

	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)

	// handle defer
	defer func() { 
		err := tp.Shutdown(ctx)
		if err != nil{
			childLogger.Error().Err(err).Msg("error closing OTEL tracer !!!")
		}
		childLogger.Info().Msg("stop done !!!")
	}()
	
	// router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(core_middleware.MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")
		json.NewEncoder(rw).Encode(appServer)
	})

	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpRouters.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpRouters.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpRouters.Header)

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/info")

		res := model.AppServer{}
		
		res.InfoPod =  appServer.InfoPod
		res.Server =  appServer.Server
		res.RedisAddress =  &appServer.DatabaseRedis.RedisAddress
		res.ConfigOTEL =  appServer.ConfigOTEL

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(res)
	})
	
	addKey := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addKey.HandleFunc("/add/key", core_middleware.MiddleWareErrorHandler(httpRouters.AddKey))		
	addKey.Use(otelmux.Middleware("go-payfee"))

	getKey := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getKey.HandleFunc("/key/{id}", core_middleware.MiddleWareErrorHandler(httpRouters.GetKey))		
	getKey.Use(otelmux.Middleware("go-payfee"))

	addScript := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addScript.HandleFunc("/add/script", core_middleware.MiddleWareErrorHandler(httpRouters.AddScript))		
	addScript.Use(otelmux.Middleware("go-payfee"))

	getScript := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getScript.HandleFunc("/script/{id}", core_middleware.MiddleWareErrorHandler(httpRouters.GetScript))		
	getScript.Use(otelmux.Middleware("go-payfee"))

	// setup http server
	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpServer.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpServer.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpServer.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpServer.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port : ", strconv.Itoa(h.httpServer.Port)).Msg("Service Port")

	// start http server
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			childLogger.Error().Err(err).Msg("canceling http mux server !!!")
		}
	}()

	// handle sigTERM
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		childLogger.Error().Err(err).Msg("warning dirty shutdown !!!")
		return
	}
}