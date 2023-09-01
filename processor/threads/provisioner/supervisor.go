package provisioner

import (
	"errors"
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/mango/api"
	"github.com/GabeCordo/mango/core/components/supervisor"
	"github.com/GabeCordo/toolchain/logging"
	"time"
)

func (thread *Thread) getSupervisor() []*supervisor.Supervisor {

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

	thread.logger.Printf("%s[%s]%s Provisioning cluster in module %s\n", logging.Green, request.Cluster, logging.Reset, request.Module)

	// Note: configs are now sent from the core, we don't need to worry about looking for, verifying, or
	//		 reverting to a default cluster.Config if one is not provided
	if request.Config == nil {
		request.Config = &clusterWrapper.DefaultConfig
	}
	supervisorInstance := clusterWrapper.CreateSupervisor(request.Metadata, thread.Config.Core, request.Config)

	thread.logger.Printf("%s[%s]%s Supervisor(%d) registered to cluster(%s)\n", logging.Green, request.Cluster, logging.Reset, supervisorInstance.Id, request.Module)

	thread.logger.Printf("%s[%s]%s Cluster Running\n", logging.Green, request.Cluster, logging.Reset)

	go func() {

		// block until the supervisor completes
		response := supervisorInstance.Start()
		// TODO : should we send the response instead?

		// TODO : define host
		api.UpdateSupervisor(thread.Config.Core, supervisorInstance.Id, supervisor.Status(supervisorInstance.State), response.Stats.ToStandard())

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
		if !clusterWrapper.IsStream() {
			thread.requestWg.Done()
		}
	}()

	// TODO : this is a temp fix for the system getting caught in an inf. deadlock
	// the issue is streams are running (listening) 24/7
	// we want to be able to shutdown the processor and unregister it from the core as an operator
	//
	// in future: discuss how to guarantee the data that is pulled is fully processed first
	// why I rushed this: this is an experimental version for NON production use
	if clusterWrapper.IsStream() {
		thread.requestWg.Done()
	}

	return nil
}
