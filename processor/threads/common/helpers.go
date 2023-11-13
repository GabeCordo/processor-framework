package common

import (
	"github.com/GabeCordo/keitt/processor/components/cluster"
	"github.com/GabeCordo/toolchain/multithreaded"
	"math/rand"
)

func SupervisorProvision(pipe chan<- ProvisionerRequest, responseTable *multithreaded.ResponseTable,
	moduleName, clusterName string, meta map[string]string, cfg *cluster.Config, timeout float64) error {

	// there is a possibility the user never passed an args value to the HTTP endpoint,
	// so we need to replace it with and empty array
	if meta == nil {
		meta = make(map[string]string)
	}
	provisionerThreadRequest := ProvisionerRequest{
		Action:   ProvisionerSupervisorCreate,
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
