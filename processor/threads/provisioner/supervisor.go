package provisioner

import (
	"errors"
	"github.com/GabeCordo/processor-framework/processor/api"
	"github.com/GabeCordo/processor-framework/processor/interfaces"
	"github.com/GabeCordo/processor-framework/processor/threads/common"
	"github.com/GabeCordo/toolchain/logging"
	"log"
	"time"
)

func (thread *Thread) getSupervisor() []*interfaces.Supervisor {

	return nil
}

func (thread *Thread) provisionSupervisor(request *common.ProvisionerRequest) error {

	// Note: all mount checks have been moved to the core

	provisionerInstance := GetProvisionerInstance()

	moduleWrapper, found := provisionerInstance.GetModule(request.Module)

	if !found {
		thread.logger.Warnf("%s[%s]%s Module does not exist\n", logging.Green, request.Module, logging.Reset)
		thread.requestWg.Done()
		return errors.New("module not found")
	}

	clusterWrapper, found := moduleWrapper.GetCluster(request.Cluster)

	if !found {
		thread.logger.Warnf("%s[%s]%s Cluster does not exist\n", logging.Green, request.Cluster, logging.Reset)
		thread.requestWg.Done()
		return errors.New("cluster not found")
	}

	// an operator shall only provision batch etl processes
	// - stream processes are meant to be run by the system when mounted or unmounted
	if (request.Source == common.User) && clusterWrapper.IsStream() {
		thread.logger.Warnf("%s[%s]%s Could not provision cluster; it's a stream process\n", logging.Green, request.Module, logging.Reset)
		thread.requestWg.Done()
		return errors.New("a stream cluster cannot be provisioned by a user")
	}

	// if we have exceeded the number of supervisors we want to run on the processor,
	// we can add it to a queue that will be run at another time.
	//
	// note: we can tell the core pre-maturely that the supervisor was provisioned so
	// 		 that the caller is told that the request successfully reached this server.
	if thread.NumOfActiveSupervisors() >= MaxNumOfSupervisors {
		if request != nil {
			thread.requestBacklog = append(thread.requestBacklog, *request)
			thread.requestWg.Done()
		} else {
			log.Printf("tried to add request to backlog but found nil pointer request")
		}
		return nil
	} else {
		thread.IncrementActiveSupervisors()
	}

	thread.logger.Printf("%s[%s]%s Provisioning cluster in module %s\n", logging.Green, request.Cluster, logging.Reset, request.Module)

	// Note: configs are now sent from the core, we don't need to worry about looking for, verifying, or
	//		 reverting to a default cluster.Config if one is not provided
	if request.Config == nil {
		request.Config = &clusterWrapper.DefaultConfig
	}
	supervisorInstance := clusterWrapper.CreateSupervisor(request.Supervisor, request.Metadata, thread.Config.Core, thread.Config.Standalone, request.Config)

	thread.logger.Printf("%s[%s]%s Supervisor(%d) registered to cluster(%s)\n", logging.Green, request.Cluster, logging.Reset, supervisorInstance.Id, request.Module)

	thread.logger.Printf("%s[%s]%s Cluster Running\n", logging.Green, request.Cluster, logging.Reset)

	go func() {

		if !thread.Config.Standalone && clusterWrapper.IsStream() {
			go func() {
				for {
					if !supervisorInstance.IsAlive() {
						break
					} else {
						api.UpdateSupervisor(thread.Config.Core, supervisorInstance.Id, interfaces.SupervisorStatus(supervisorInstance.State), supervisorInstance.Stats.ToStandard())
					}

					time.Sleep(1 * time.Second)
				}
			}()
		}

		// block until the supervisor completes
		response := supervisorInstance.Start()
		// TODO : should we send the response instead?

		// TODO : define host
		if !thread.Config.Standalone {
			api.UpdateSupervisor(thread.Config.Core, supervisorInstance.Id, interfaces.SupervisorStatus(supervisorInstance.State), response.Stats.ToStandard())
		}

		// provide the console with output indicating that the cluster has completed
		// we already provide output when a cluster is provisioned, so it completes the state
		if thread.Config.Debug {
			duration := time.Now().Sub(supervisorInstance.StartTime)
			thread.logger.Printf("%s[%s]%s Cluster transformations complete, took %dhr %dm %ds %dms %dus\n",
				logging.Green,
				supervisorInstance.Config.Identifier,
				logging.Reset,
				int(duration.Hours()),
				int(duration.Minutes()),
				int(duration.Seconds()),
				int(duration.Milliseconds()),
				int(duration.Microseconds()),
			)
		}

		// let the provisioner thread decrement the semaphore otherwise we will be stuck in deadlock waiting for
		// the provisioned cluster to complete before allowing the etl-threads to shut down
		//if !clusterWrapper.IsStream() {
		thread.DecrementActiveSupervisors()
		thread.requestWg.Done()
		//}
	}()

	// TODO : this is a temp fix for the system getting caught in an inf. deadlock
	// the issue is streams are running (listening) 24/7
	// we want to be able to shutdown the processor and unregister it from the core as an operator
	//
	// in future: discuss how to guarantee the data that is pulled is fully processed first
	// why I rushed this: this is an experimental version for NON production use
	//if clusterWrapper.IsStream() {
	//	thread.requestWg.Done()
	//}

	return nil
}
