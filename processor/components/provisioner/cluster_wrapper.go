package provisioner

import (
	"github.com/GabeCordo/processor-framework/processor/components/cluster"
	"github.com/GabeCordo/processor-framework/processor/components/supervisor"
)

func NewClusterWrapper(moduleName, identifier string, mode cluster.EtlMode, implementation cluster.Cluster) *ClusterWrapper {

	clusterWrapper := new(ClusterWrapper)

	clusterWrapper.registry = supervisor.NewRegistry(moduleName, identifier, implementation)
	clusterWrapper.Identifier = identifier
	clusterWrapper.Module = moduleName
	clusterWrapper.Mode = mode
	clusterWrapper.Mounted = false

	return clusterWrapper
}

func (clusterWrapper *ClusterWrapper) IsStream() bool {

	return clusterWrapper.Mode == cluster.Stream
}

func (clusterWrapper *ClusterWrapper) IsMounted() bool {

	return clusterWrapper.Mounted
}

func (clusterWrapper *ClusterWrapper) Mount() *ClusterWrapper {

	clusterWrapper.Mounted = true
	return clusterWrapper
}

func (clusterWrapper *ClusterWrapper) UnMount() *ClusterWrapper {

	clusterWrapper.Mounted = false
	return clusterWrapper
}

func (clusterWrapper *ClusterWrapper) GetClusterImplementation() cluster.Cluster {
	return clusterWrapper.registry.GetClusterImplementation()
}

func (clusterWrapper *ClusterWrapper) FindSupervisors() []*supervisor.Supervisor {
	return clusterWrapper.registry.GetSupervisors()
}

func (clusterWrapper *ClusterWrapper) FindSupervisor(id uint64) (instance *supervisor.Supervisor, found bool) {

	instance, found = clusterWrapper.registry.GetSupervisor(id)
	return instance, found
}

func (clusterWrapper *ClusterWrapper) CreateSupervisor(metadata map[string]string, core string, standalone bool, config ...*cluster.Config) *supervisor.Supervisor {

	return clusterWrapper.registry.CreateSupervisor(metadata, core, standalone, config...)
}

func (clusterWrapper *ClusterWrapper) DeleteSupervisor(identifier uint64) (deleted, found bool) {

	deleted, found = clusterWrapper.registry.DeleteSupervisor(identifier)
	return deleted, found
}

func (clusterWrapper *ClusterWrapper) SuspendSupervisors() {

	clusterWrapper.registry.SuspendSupervisors()
}

func (clusterWrapper *ClusterWrapper) CanDelete() (canDelete bool) {

	canDelete = true
	for _, supervisorInstance := range clusterWrapper.registry.GetSupervisors() {
		if !supervisorInstance.Deletable() {
			canDelete = false
			break
		}
	}

	return canDelete
}
