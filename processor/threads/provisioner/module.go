package provisioner

import (
	"errors"
	"github.com/GabeCordo/mango-go/processor/components/module"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/proxy"
	"github.com/GabeCordo/mango/utils"
)

func (thread *Thread) ProcessAddModule(request *common.ProvisionerRequest) error {

	defer thread.requestWg.Done()

	thread.logger.Printf("registering module at %s\n", request.Path)

	remoteModule, err := module.NewRemoteModule(request.Path)
	if err != nil {
		thread.logger.Alertln("cannot find remote module")
		return errors.New("cannot find remote module")
	}

	moduleInstance, err := remoteModule.Get()
	if err != nil {
		thread.logger.Alertln("module built with older version")
		return errors.New("module built with older version")
	}

	if err := GetProvisionerInstance().AddModule(moduleInstance); err != nil {
		thread.logger.Alertln("a module with that identifier already exists or is corrupt")
		return errors.New("a module with that identifier already exists or is corrupt")
	}

	// note: the module wrapper should already be defined so there is no need to validate
	moduleWrapper, _ := GetProvisionerInstance().GetModule(moduleInstance.Config.Name)

	registeredClusters := moduleWrapper.GetClusters()

	// REGISTER ANY HELPERS TO CLUSTERS THAT HAVE DEFINED THEM WITHIN THE STRUCT
	for _, clusterWrapper := range registeredClusters {

		clusterImplementation := clusterWrapper.GetClusterImplementation()

		if helperImplementation, ok := (clusterImplementation).(utils.UsesHelper); ok {
			helper, _ := utils.NewStandardHelper(moduleWrapper.Identifier, clusterWrapper.Identifier)
			helperImplementation.SetHelper(helper)
		}
	}

	err = proxy.CreateModule(common.GetConfigInstance().Core, moduleInstance.Config)
	return err
}

func (thread *Thread) ProcessDeleteModule(request *common.ProvisionerRequest) error {

	defer thread.requestWg.Done()

	provisionerInstance := GetProvisionerInstance()

	//response := &threads.ProvisionerResponse{Nonce: request.Nonce}
	deleted, _, found := provisionerInstance.DeleteModule(request.Module)

	//response.Success = true
	//if deleted {
	//	response.Description = "module deleted"
	//	//
	//	//databaseRequest := &threads.DatabaseRequest{
	//	//	Action: threads.DatabaseDelete,
	//	//	Type:   threads.ClusterModule,
	//	//	Module: request.ModuleName,
	//	//	Nonce:  rand.Uint32(),
	//	//}
	//	//thread.Request(threads.Database, databaseRequest)
	//	//
	//	//data, didTimeout := utils.SendAndWait(thread.databaseResponseTable, databaseRequest.Nonce,
	//	//	common.GetConfigInstance().MaxWaitForResponse)
	//	//
	//	//if didTimeout {
	//	//	response.Success = false
	//	//	response.Description = "could not delete clusters and statistics under a module"
	//	//}
	//	//
	//	//databaseResponse := (data).(threads.DatabaseResponse)
	//	//
	//	//if !databaseResponse.Success {
	//	//	response.Success = false
	//	//	response.Description = "could not delete clusters and statistics under a module"
	//	//}
	//} else {
	//	response.Description = "module marked for deletion, a cluster is likely running right now, try later"
	//}
	//} else {
	//	response.Success = false
	//	response.Description = "module not found"
	//}

	if !found {
		return errors.New("module does not exist on the processor")
	}

	if deleted {
		thread.logger.Printf("locally the module (%s) has been deleted\n", request.Module)
	} else {
		thread.logger.Printf("locally the module (%s) was marked for deletion\n", request.Module)
	}

	err := proxy.DeleteModule(common.GetConfigInstance().Core, request.Module)

	if err != nil {
		return errors.New("it's likely the core didn't delete the module")
	}

	return nil
}