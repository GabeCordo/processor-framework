package provisioner

import (
	"errors"
	"fmt"
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/mango/api"
)

func (thread *Thread) Setup() {

	thread.accepting = true

	// initialize a provisioner instance with a common module
	GetProvisionerInstance()
}

func (thread *Thread) Start() {

	// INCOMING REQUESTS

	go func() {
		for request := range thread.C1 {
			if !thread.accepting {
				break
			}
			thread.requestWg.Add(1)
			thread.processRequest(&request)
		}

		thread.listenersWg.Wait()
	}()

	// STANDALONE MODE

	if thread.Config.Standalone {

		for _, moduleInst := range GetProvisionerInstance().GetModules() {

			for _, clusterInst := range moduleInst.GetClusters() {

				if !clusterInst.IsStream() {
					continue
				}

				request := &common.ProvisionerRequest{
					Action:   common.ProvisionerSupervisorCreate,
					Source:   common.Core,
					Module:   moduleInst.Identifier,
					Cluster:  clusterInst.Identifier,
					Config:   &clusterInst.DefaultConfig,
					Metadata: make(map[string]string),
				}
				thread.requestWg.Add(1)
				thread.provisionSupervisor(request)
			}
		}
	} else {

		for _, moduleInst := range GetProvisionerInstance().GetModules() {
			api.CreateModule(thread.Config.Core, &thread.Config.Processor, moduleInst.ToConfig())
		}
	}

	thread.listenersWg.Wait()
	thread.requestWg.Wait()
}

func (thread *Thread) respond(response *common.ProvisionerResponse) {

	thread.C2 <- *response
}

func (thread *Thread) processRequest(request *common.ProvisionerRequest) {

	response := &common.ProvisionerResponse{Error: nil, Nonce: request.Nonce}

	switch request.Action {
	case common.ProvisionerModuleGet:
		response.Data = thread.getModules()
	case common.ProvisionerSupervisorCreate:
		response.Error = thread.provisionSupervisor(request)
	default:
		response.Error = errors.New("bad request")
		thread.requestWg.Done()
	}

	response.Success = response.Error == nil
	thread.respond(response)
}

func (thread *Thread) Teardown() {
	thread.accepting = false

	modules := GetProvisionerInstance().GetModules()

	for _, module := range modules {

		clusters := module.GetClusters()

		for _, cluster := range clusters {

			if !cluster.IsStream() {
				continue
			}

			supervisors := cluster.FindSupervisors()

			for _, supervisor := range supervisors {

				if !supervisor.IsAlive() {

					fmt.Printf("supervisor is not alive %d %s\n", supervisor.Id, supervisor.State.ToString())
					continue
				}

				fmt.Printf("marking supervisor as teardown %d\n", supervisor.Id)
				supervisor.Teardown()
			}
		}
	}

	thread.requestWg.Wait()
}
