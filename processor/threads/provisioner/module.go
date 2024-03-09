package provisioner

import "github.com/GabeCordo/processor-framework/processor/components/provisioner"

func (thread *Thread) getModules() []*provisioner.ModuleWrapper {

	defer thread.requestWg.Done()
	return GetProvisionerInstance().GetModules()
}
