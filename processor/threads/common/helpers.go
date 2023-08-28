package common

import (
	"github.com/GabeCordo/keitt/processor/components/cluster"
	"github.com/GabeCordo/toolchain/multithreaded"
	"math/rand"
)

func RegisterModule(pipe chan<- ProvisionerRequest, responseTable *multithreaded.ResponseTable,
	modulePath string, timeout float64) error {

	request := ProvisionerRequest{
		Action: ProvisionerCreateModule,
		Path:   modulePath,
		Nonce:  rand.Uint32(),
	}
	pipe <- request

	rsp, didTimeout := multithreaded.SendAndWait(responseTable, request.Nonce, timeout)
	if didTimeout {
		return multithreaded.NoResponseReceived
	}

	response := (rsp).(ProvisionerResponse)
	return response.Error
}

func DeleteModule(pipe chan<- ProvisionerRequest, responseTable *multithreaded.ResponseTable,
	moduleName string, timeout float64) error {

	request := ProvisionerRequest{
		Action: ProvisionerDeleteModule,
		Module: moduleName,
		Nonce:  rand.Uint32(),
	}
	pipe <- request

	data, didTimeout := multithreaded.SendAndWait(responseTable, request.Nonce, timeout)
	if didTimeout {
		return multithreaded.NoResponseReceived
	}

	response := (data).(ProvisionerResponse)
	return response.Error
}

func SupervisorProvision(pipe chan<- ProvisionerRequest, responseTable *multithreaded.ResponseTable,
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

	data, didTimeout := multithreaded.SendAndWait(responseTable, provisionerThreadRequest.Nonce, timeout)
	if didTimeout {
		return multithreaded.NoResponseReceived
	}

	provisionerResponse := (data).(ProvisionerResponse)
	return provisionerResponse.Error
}

func ShutdownCore(pipe chan<- InterruptEvent) {
	pipe <- Shutdown
}
