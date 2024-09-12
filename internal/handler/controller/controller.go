package controller

import (	
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"

	"github.com/go-payfee/internal/service"
	"github.com/go-payfee/internal/core"
	"github.com/go-payfee/internal/erro"
	"github.com/go-payfee/internal/lib"
)

var childLogger = log.With().Str("handler", "controller").Logger()

type APIError struct {
	StatusCode	int  `json:"statusCode"`
	Msg			any `json:"msg"`
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
	}
}

type HttpWorkerAdapter struct {
	workerService 	*service.RedisService
}

func NewHttpWorkerAdapter(workerService *service.RedisService) HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	
	return HttpWorkerAdapter{
		workerService: workerService,
	}
}

func (h *HttpWorkerAdapter) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
	return
}

func (h *HttpWorkerAdapter) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
	return
}

func (h *HttpWorkerAdapter) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Header")
	
	json.NewEncoder(rw).Encode(req.Header)
	return
}

func (h *HttpWorkerAdapter) GetScript(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetScript")

	span := lib.Span(req.Context(), "handler.GetScript")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	script := core.ScriptData{}
	script.Script.Name = varID
	
	res, err := h.workerService.GetScript(req.Context(), script)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(404, err)
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	//rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) AddScript( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("AddScript")

	span := lib.Span(req.Context(), "handler.AddScript")
	defer span.End()

	script := core.ScriptData{}
	err := json.NewDecoder(req.Body).Decode(&script)
    if err != nil {
		apiError := NewAPIError(400, erro.ErrUnmarshal)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    }

	res, err := h.workerService.AddScript(req.Context(), script)
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetKey(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetKey")

	span := lib.Span(req.Context(), "handler.GetKey")
	defer span.End()

	vars := mux.Vars(req)
	varID := vars["id"]

	fee := core.Fee{}
	fee.Name = varID
	
	res, err := h.workerService.GetKey(req.Context(), fee)
	if err != nil {
		var apiError APIError
		switch err {
			case erro.ErrNotFound:
				apiError = NewAPIError(404, err)
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) AddKey( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("SetKey")

	span := lib.Span(req.Context(), "handler.SetKey")
	defer span.End()

	fee := core.Fee{}
	err := json.NewDecoder(req.Body).Decode(&fee)
    if err != nil {
		apiError := NewAPIError(400, erro.ErrUnmarshal)
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
    }

	res, err := h.workerService.AddKey(req.Context(), fee)
	if err != nil {
		var apiError APIError
		switch err {
			default:
				apiError = NewAPIError(500, err)
		}
		rw.WriteHeader(apiError.StatusCode)
		json.NewEncoder(rw).Encode(apiError)
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}