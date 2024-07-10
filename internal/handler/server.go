package handler

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"context"

	"github.com/gorilla/mux"

	"github.com/go-payfee/internal/service"
	"github.com/go-payfee/internal/core"
	"github.com/go-payfee/internal/lib"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

//----------------------------------------------------------------
type HttpWorkerAdapter struct {
	workerService 	*service.RedisService
}

func NewHttpWorkerAdapter(workerService *service.RedisService) HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	
	return HttpWorkerAdapter{
		workerService: workerService,
	}
}
//------------------------------------------------------
type HttpServer struct {
	httpServer	*core.Server
}

func NewHttpAppServer(httpServer *core.Server) HttpServer {
	childLogger.Debug().Msg("NewHttpAppServer")

	return HttpServer{httpServer: httpServer }
}
//-------------------------------------------------
func (h HttpServer) StartHttpAppServer(	ctx context.Context, 
										httpWorkerAdapter *HttpWorkerAdapter,
										appServer *core.AppServer) {
	childLogger.Info().Msg("StartHttpAppServer")
	// ---------------------- OTEL ---------------
	childLogger.Info().Str("OTEL_EXPORTER_OTLP_ENDPOINT :", appServer.ConfigOTEL.OtelExportEndpoint).Msg("")
	
	tp := lib.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoPod)
	defer func() { 
		err := tp.Shutdown(ctx)
		if err != nil{
			childLogger.Error().Err(err).Msg("Erro closing OTEL tracer !!!")
		}
	}()
	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")

		json.NewEncoder(rw).Encode(appServer)
	})

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/info")
		
		res := core.AppServer{}
		
		res.InfoPod =  appServer.InfoPod
		res.Server =  appServer.Server
		res.RedisAddress =  appServer.RedisAddress
		res.ConfigOTEL =  appServer.ConfigOTEL

		json.NewEncoder(rw).Encode(res)
	})
	
	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpWorkerAdapter.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpWorkerAdapter.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpWorkerAdapter.Header)
	header.Use(otelmux.Middleware("go-payfee"))

	setKey := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	setKey.Handle("/key/add", 
						http.HandlerFunc(httpWorkerAdapter.AddKey),)
	setKey.Use(otelmux.Middleware("go-payfee"))

	getKey := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getKey.Handle("/key/get/{id}",
						http.HandlerFunc(httpWorkerAdapter.GetKey),)
	getKey.Use(otelmux.Middleware("go-payfee"))

	addScript := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addScript.Handle("/script/add", 
						http.HandlerFunc(httpWorkerAdapter.AddScript),)
	addScript.Use(otelmux.Middleware("go-payfee"))

	getScript := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getScript.Handle("/script/get/{id}",
						http.HandlerFunc(httpWorkerAdapter.GetScript),)
						getScript.Use(otelmux.Middleware("go-payfee"))

	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpServer.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpServer.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpServer.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpServer.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port : ", strconv.Itoa(h.httpServer.Port)).Msg("Service Port")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			childLogger.Error().Err(err).Msg("Cancel http mux server !!!")
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		childLogger.Error().Err(err).Msg("WARNING Dirty Shutdown !!!")
		return
	}
	childLogger.Info().Msg("Stop Done !!!!")
}