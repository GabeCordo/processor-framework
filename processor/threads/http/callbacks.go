package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/components/cluster"
	"github.com/GabeCordo/mango/core"
	"github.com/GabeCordo/mango/utils"
	"net/http"
	"time"
)

type ModuleRequestBody struct {
	ModulePath string `json:"path,omitempty"`
	ModuleName string `json:"module,omitempty"`
}

func (httpThread *Thread) moduleCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		httpThread.postModuleCallback(w, r)
	} else if r.Method == "DELETE" {
		httpThread.deleteModuleCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (httpThread *Thread) postModuleCallback(w http.ResponseWriter, r *http.Request) {

	request := &ModuleRequestBody{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.RegisterModule(httpThread.C1, httpThread.ProvisionerResponseTable, request.ModulePath)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusConflict)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

func (httpThread *Thread) deleteModuleCallback(w http.ResponseWriter, r *http.Request) {

	request := &ModuleRequestBody{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.DeleteModule(httpThread.C1, httpThread.ProvisionerResponseTable, request.ModuleName)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

type JSONResponse struct {
	Status      int    `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	Data        any    `json:"data,omitempty"`
}

type SupervisorConfigJSONBody struct {
	Module     string            `json:"module"`
	Cluster    string            `json:"cluster"`
	Config     cluster.Config    `json:"common"`
	Supervisor uint64            `json:"id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type SupervisorProvisionJSONResponse struct {
	Cluster    string `json:"cluster,omitempty"`
	Supervisor uint64 `json:"id,omitempty"`
}

func (httpThread *Thread) supervisorCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		httpThread.postSupervisorCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (httpThread *Thread) postSupervisorCallback(w http.ResponseWriter, r *http.Request) {

	var request SupervisorConfigJSONBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if (r.Method != "GET") && (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.SupervisorProvision(httpThread.C1, httpThread.ProvisionerResponseTable,
		request.Module, request.Cluster, request.Metadata, &request.Config)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

type DebugJSONBody struct {
	Action string `json:"action"`
}

type DebugJSONResponse struct {
	Duration time.Duration `json:"time-elapsed"`
	Success  bool          `json:"success"`
}

func (httpThread *Thread) debugCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		httpThread.getDebugCallback(w, r)
	} else if r.Method == "POST" {
		httpThread.postDebugCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (httpThread *Thread) getDebugCallback(w http.ResponseWriter, r *http.Request) {
	// do nothing
}

func (httpThread *Thread) postDebugCallback(w http.ResponseWriter, r *http.Request) {

	request := &DebugJSONBody{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if (r.Method != "OPTIONS") && err != nil {
		fmt.Println("missing body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := core.Response{Success: true}

	if request.Action == "shutdown" {
		common.ShutdownCore(httpThread.Interrupt)
	} else if request.Action == "debug" {
		response.Description = common.ToggleDebugMode(httpThread.logger)
	}

	b, _ := json.Marshal(response)
	w.Write(b)
}
