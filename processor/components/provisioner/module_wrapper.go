package provisioner

import (
	"errors"
	"fmt"
	"github.com/GabeCordo/processor-framework/processor/components/cluster"
	"github.com/GabeCordo/processor-framework/processor/interfaces"
	"log"
)

func NewModuleWrapper() *ModuleWrapper {

	moduleWrapper := new(ModuleWrapper)

	moduleWrapper.clusters = make(map[string]*ClusterWrapper)
	moduleWrapper.Mounted = false

	return moduleWrapper
}

func (moduleWrapper *ModuleWrapper) IsMounted() bool {

	return moduleWrapper.Mounted
}

func (moduleWrapper *ModuleWrapper) Mount() *ModuleWrapper {

	moduleWrapper.Mounted = true
	return moduleWrapper
}

func (moduleWrapper *ModuleWrapper) UnMount() *ModuleWrapper {

	moduleWrapper.Mounted = false
	return moduleWrapper
}

func (moduleWrapper *ModuleWrapper) GetClustersData() map[string]bool {

	mounts := make(map[string]bool)

	for identifier, clusterWrapper := range moduleWrapper.clusters {
		mounts[identifier] = clusterWrapper.Mounted
	}

	return mounts
}

func (moduleWrapper *ModuleWrapper) GetClusters() (clusterWrappers []*ClusterWrapper) {

	clusterWrappers = make([]*ClusterWrapper, 0)

	for _, clusterWrapper := range moduleWrapper.clusters {
		clusterWrappers = append(clusterWrappers, clusterWrapper)
	}

	return clusterWrappers
}

func (moduleWrapper *ModuleWrapper) GetCluster(clusterName string) (clusterWrapper *ClusterWrapper, found bool) {

	moduleWrapper.mutex.RLock()
	defer moduleWrapper.mutex.RUnlock()

	clusterWrapper, found = moduleWrapper.clusters[clusterName]
	return clusterWrapper, found
}

// AddCluster
// Creates a new cluster record within the calling module. The cluster name defines the keyword
// an operator uses to provision a cluster, and the mode represents how the cluster is run.
func (moduleWrapper *ModuleWrapper) AddCluster(clusterName string, mode string, implementation cluster.Cluster, cfg ...*cluster.Config) (*ClusterWrapper, error) {

	moduleWrapper.mutex.RLock()

	if _, found := moduleWrapper.clusters[clusterName]; found {
		return nil, errors.New("a cluster with this identifier already exists")
	}

	moduleWrapper.mutex.RUnlock()

	moduleWrapper.mutex.Lock()
	defer moduleWrapper.mutex.Unlock()

	clusterWrapper := NewClusterWrapper(moduleWrapper.Identifier, clusterName, cluster.EtlMode(mode), implementation)
	clusterWrapper.Mounted = true
	clusterWrapper.DefaultConfig = cluster.DefaultConfig  // copy
	clusterWrapper.DefaultConfig.Identifier = clusterName // copy

	for _, c := range cfg {
		clusterWrapper.DefaultConfig = *c
		clusterWrapper.DefaultConfig.Identifier = clusterName
	}

	moduleWrapper.clusters[clusterName] = clusterWrapper

	return clusterWrapper, nil
}

func (moduleWrapper *ModuleWrapper) DeleteCluster(identifier string) (deleted, found bool) {

	clusterWrapper, found := moduleWrapper.clusters[identifier]
	if !found {
		return false, false
	}

	if !clusterWrapper.CanDelete() {
		return false, true
	}

	moduleWrapper.mutex.Lock()
	defer moduleWrapper.mutex.Unlock()

	delete(moduleWrapper.clusters, identifier)
	return true, true
}

func (moduleWrapper *ModuleWrapper) CanDelete() (canDelete bool) {

	moduleWrapper.mutex.RLock()
	defer moduleWrapper.mutex.RUnlock()

	// if the module is not marked for deletion, it should not be deleted
	if !moduleWrapper.MarkForDeletion {
		log.Printf("[provisioner] cannot delete %s - not marked for deletion\n", moduleWrapper.Identifier)
		return false
	}

	canDelete = true
	// look over all the supervisor in a module
	for clusterName, clusterWrapper := range moduleWrapper.clusters {

		if !clusterWrapper.CanDelete() {
			log.Printf("[provisioner][cluster] cannot delete %s\n", clusterName)
			canDelete = false
			break
		}
	}

	return canDelete
}

func (moduleWrapper *ModuleWrapper) ToConfig() *interfaces.ModuleConfig {
	cfg := new(interfaces.ModuleConfig)

	cfg.Name = moduleWrapper.Identifier
	cfg.Version = moduleWrapper.Version
	cfg.Exports = make([]interfaces.Cluster, len(moduleWrapper.clusters))

	idx := 0
	for _, cluster := range moduleWrapper.clusters {
		cfg.Exports[idx] = interfaces.Cluster{
			Cluster:     cluster.Identifier,
			StaticMount: cluster.Mounted,
			Config: interfaces.ClusterConfig{
				Mode:    interfaces.EtlMode(cluster.Mode),
				OnCrash: interfaces.OnCrash(cluster.DefaultConfig.OnCrash),
				OnLoad:  interfaces.OnLoad(cluster.DefaultConfig.OnLoad),
				Static: struct {
					TFunctions int `yaml:"t-functions" json:"t-functions"`
					LFunctions int `yaml:"l-functions" json:"l-functions"`
				}{
					TFunctions: cluster.DefaultConfig.StartWithNTransformClusters,
					LFunctions: cluster.DefaultConfig.StartWithNTransformClusters,
				},
				Dynamic: struct {
					TFunction interfaces.DynamicFeatures `yaml:"t-function" json:"t-function"`
					LFunction interfaces.DynamicFeatures `yaml:"l-function" json:"l-function"`
				}{
					TFunction: interfaces.DynamicFeatures{
						Threshold:    cluster.DefaultConfig.ETChannelThreshold,
						GrowthFactor: cluster.DefaultConfig.ETChannelGrowthFactor,
					},
					LFunction: interfaces.DynamicFeatures{
						Threshold:    cluster.DefaultConfig.TLChannelThreshold,
						GrowthFactor: cluster.DefaultConfig.TLChannelGrowthFactor,
					},
				},
			},
		}
		idx++
	}

	return cfg
}

func (moduleWrapper *ModuleWrapper) Print() {
	fmt.Printf("%s %.3f\n", moduleWrapper.Identifier, moduleWrapper.Version)
}
