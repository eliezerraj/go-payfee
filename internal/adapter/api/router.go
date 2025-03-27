package api

import (
	"encoding/json"
	"net/http"
	"github.com/rs/zerolog/log"

	"github.com/go-payfee/internal/core/service"
	"github.com/go-payfee/internal/core/model"
	"github.com/go-payfee/internal/core/erro"
	go_core_observ "github.com/eliezerraj/go-core/observability"
	"github.com/eliezerraj/go-core/coreJson"
	"github.com/gorilla/mux"
)

var childLogger = log.With().Str("component", "go-payfee").Str("package", "internal.adapter.api").Logger()

var core_json coreJson.CoreJson
var core_apiError coreJson.APIError
var tracerProvider go_core_observ.TracerProvider

type HttpRouters struct {
	workerService 	*service.WorkerService
}

func NewHttpRouters(workerService *service.WorkerService) HttpRouters {
	childLogger.Info().Str("func","NewHttpRouters").Send()
	
	return HttpRouters{
		workerService: workerService,
	}
}

// About return a health
func (h *HttpRouters) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Health").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About return a live
func (h *HttpRouters) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Live").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	json.NewEncoder(rw).Encode(model.MessageRouter{Message: "true"})
}

// About show all header received
func (h *HttpRouters) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Info().Str("func","Header").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()
	
	json.NewEncoder(rw).Encode(req.Header)
}

// About get all script
func (h *HttpRouters) GetScript(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","GetScript").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.GetScript")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	script := model.ScriptData{}
	script.Script.Name = varID

	res, err := h.workerService.GetScript(req.Context(), &script)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About add script
func (h *HttpRouters) AddScript(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","AddScript").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.AddScript")
	defer span.End()

	script := model.ScriptData{}
	err := json.NewDecoder(req.Body).Decode(&script)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.AddScript(req.Context(), &script)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About get key
func (h *HttpRouters) GetKey(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","GetKey").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.GetKey")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	fee := model.Fee{}
	fee.Name = varID

	res, err := h.workerService.GetKey(req.Context(), &fee)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}

// About add key
func (h *HttpRouters) AddKey(rw http.ResponseWriter, req *http.Request) error {
	childLogger.Info().Str("func","AddKey").Interface("trace-resquest-id", req.Context().Value("trace-request-id")).Send()

	span := tracerProvider.Span(req.Context(), "adapter.api.AddKey")
	defer span.End()

	fee := model.Fee{}
	err := json.NewDecoder(req.Body).Decode(&fee)
    if err != nil {
		core_apiError = core_apiError.NewAPIError(err, http.StatusBadRequest)
		return &core_apiError
    }
	defer req.Body.Close()

	res, err := h.workerService.AddKey(req.Context(), &fee)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			core_apiError = core_apiError.NewAPIError(err, http.StatusNotFound)
		default:
			core_apiError = core_apiError.NewAPIError(err, http.StatusInternalServerError)
		}
		return &core_apiError
	}
	
	return core_json.WriteJSON(rw, http.StatusOK, res)
}