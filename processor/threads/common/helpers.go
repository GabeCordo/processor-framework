package common

import (
	"github.com/GabeCordo/mango/components/cluster"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
	"math/rand"
)

func RegisterModule(pipe chan<- ProvisionerRequest, responseTable *utils.ResponseTable,
	modulePath string, timeout float64) error {

	request := ProvisionerRequest{
		Action: ProvisionerCreateModule,
		Path:   modulePath,
		Nonce:  rand.Uint32(),
	}
	pipe <- request

	rsp, didTimeout := utils.SendAndWait(responseTable, request.Nonce, timeout)
	if didTimeout {
		return utils.NoResponseReceived
	}

	response := (rsp).(ProvisionerResponse)
	return response.Error
}

func DeleteModule(pipe chan<- ProvisionerRequest, responseTable *utils.ResponseTable,
	moduleName string, timeout float64) error {

	request := ProvisionerRequest{
		Action: ProvisionerDeleteModule,
		Module: moduleName,
		Nonce:  rand.Uint32(),
	}
	pipe <- request

	data, didTimeout := utils.SendAndWait(responseTable, request.Nonce, timeout)
	if didTimeout {
		return utils.NoResponseReceived
	}

	response := (data).(ProvisionerResponse)
	return response.Error
}

func SupervisorProvision(pipe chan<- ProvisionerRequest, responseTable *utils.ResponseTable,
	moduleName, clusterName string, meta map[string]string, cfg *cluster.Config, timeout float64) error {

	// there is a possibility the user never passed an args value to the HTTP endpoint,
	// so we need to replace it with and empty arry
	if meta == nil {
		meta = make(map[string]string)
	}
	provisionerThreadRequest := ProvisionerRequest{
		Module:   moduleName,
		Cluster:  clusterName,
		Metadata: meta,
		Config:   cfg,
		Nonce:    rand.Uint32(),
	}
	pipe <- provisionerThreadRequest

	data, didTimeout := utils.SendAndWait(responseTable, provisionerThreadRequest.Nonce, timeout)
	if didTimeout {
		return utils.NoResponseReceived
	}

	provisionerResponse := (data).(ProvisionerResponse)
	return provisionerResponse.Error
}

func ShutdownCore(pipe chan<- threads.InterruptEvent) {
	pipe <- threads.Shutdown
}
