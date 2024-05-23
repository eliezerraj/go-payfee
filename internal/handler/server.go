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
	"fmt"

	"github.com/gorilla/mux"

	"github.com/go-payfee/internal/service"
	"github.com/go-payfee/internal/core"
	"github.com/aws/aws-xray-sdk-go/xray"

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

	setKey := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	setKey.Handle("/key/add", 
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "payfee:", appServer.InfoPod.AvailabilityZone, "./key/add")), 
						http.HandlerFunc(httpWorkerAdapter.AddKey),
						),
	)

	getKey := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getKey.Handle("/key/get/{id}",
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "payfee:", appServer.InfoPod.AvailabilityZone, "./key/get")),
						http.HandlerFunc(httpWorkerAdapter.GetKey),
						),
	)
	
	addScript := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addScript.Handle("/script/add", 
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "payfee:", appServer.InfoPod.AvailabilityZone, "./script/add")), 
						http.HandlerFunc(httpWorkerAdapter.AddScript),
						),
	)

	getScript := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getScript.Handle("/script/get/{id}",
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "payfee:", appServer.InfoPod.AvailabilityZone, "./script/get")),
						http.HandlerFunc(httpWorkerAdapter.GetScript),
						),
	)

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