package provisioner

import (
	"errors"
	"github.com/GabeCordo/keitt/processor/threads/common"
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
			if !moduleInst.IsMounted() {
				continue
			}

			for _, clusterInst := range moduleInst.GetClusters() {
				if !clusterInst.IsMounted() {
					continue
				}

				if !clusterInst.IsStream() {
					continue
				}

				request := &common.ProvisionerRequest{
					Action:   common.ProvisionerCreateSupervisor,
					Source:   common.Core,
					Module:   moduleInst.Identifier,
					Cluster:  clusterInst.Identifier,
					Config:   &clusterInst.DefaultConfig,
					Metadata: make(map[string]string),
				}
				thread.requestWg.Add(1)
				thread.ProcessProvisionRequest(request)
			}
		}
	}

	thread.listenersWg.Wait()
	thread.requestWg.Wait()
}

func (thread *Thread) respond(response *common.ProvisionerResponse) {

	thread.C2 <- *response
}

func (thread *Thread) processRequest(request *common.ProvisionerRequest) {

	response := &common.ProvisionerResponse{Error: nil}

	switch request.Action {
	case common.ProvisionerCreateSupervisor:
		response.Error = thread.ProcessProvisionRequest(request)
	case common.ProvisionerCreateModule:
		response.Error = thread.ProcessAddModule(request)
	case common.ProvisionerDeleteModule:
		response.Error = thread.ProcessDeleteModule(request)
	default:
		response.Error = errors.New("bad request")
	}

	response.Success = response.Error == nil
	thread.respond(response)
}

func (thread *Thread) Teardown() {
	thread.accepting = false

	thread.requestWg.Wait()
}
