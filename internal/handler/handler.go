package handler

import (	
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"

	"github.com/go-payfee/internal/core"
	"github.com/go-payfee/internal/erro"
	
)

var childLogger = log.With().Str("handler", "handler").Logger()

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

	vars := mux.Vars(req)
	varID := vars["id"]

	script := core.ScriptData{}
	script.Script.Name = varID
	
	res, err := h.workerService.GetScript(req.Context(), script)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			rw.WriteHeader(404)
			json.NewEncoder(rw).Encode(err.Error())
			return
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) AddScript( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("AddScript")

	script := core.ScriptData{}
	err := json.NewDecoder(req.Body).Decode(&script)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.AddScript(req.Context(), script)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetKey(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetKey")

	vars := mux.Vars(req)
	varID := vars["id"]

	fee := core.Fee{}
	fee.Name = varID
	
	res, err := h.workerService.GetKey(req.Context(), fee)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			rw.WriteHeader(404)
			json.NewEncoder(rw).Encode(err.Error())
			return
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) AddKey( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("SetKey")

	fee := core.Fee{}
	err := json.NewDecoder(req.Body).Decode(&fee)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.AddKey(req.Context(), fee)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}