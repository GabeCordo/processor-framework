package supervisor

import (
	"github.com/GabeCordo/processor-framework/processor/components/channel"
	"github.com/GabeCordo/processor-framework/processor/components/cluster"
	"sync"
	"time"
)

const (
	MaxConcurrentSupervisors = 24
)

type Status string

const (
	UnTouched    Status = "untouched"
	Running             = "running"
	Provisioning        = "provisioning"
	Failed              = "failed"
	Stopping            = "stopping"
	Terminated          = "terminated"
	Unknown             = "-"
)

type Event uint8

const (
	Startup Event = iota
	StartProvision
	EndProvision
	Error
	Suspend
	TearedDown
	StartReport
	EndReport
)

type Supervisor struct {
	Id uint64 `json:"id"`

	Config    cluster.Config      `json:"common"`
	Stats     *cluster.Statistics `json:"stats"`
	State     Status              `json:"status"`
	Mode      cluster.OnCrash     `json:"on-crash"`
	StartTime time.Time           `json:"start-time"`

	Metadata cluster.M `json:"meta-data"`

	group     cluster.Cluster
	helper    cluster.H
	ETChannel *channel.ManagedChannel
	TLChannel *channel.ManagedChannel

	loadWaitGroup sync.WaitGroup
	waitGroup     sync.WaitGroup
	mutex         sync.RWMutex
}

func NewSupervisor(clusterImplementation cluster.Cluster, metadata map[string]string, helper cluster.H) *Supervisor {
	supervisor := new(Supervisor)

	supervisor.group = clusterImplementation
	supervisor.State = UnTouched
	supervisor.Config = cluster.DefaultConfig
	supervisor.Stats = cluster.NewStatistics()
	supervisor.ETChannel = channel.NewManagedChannel("ETChannel", supervisor.Config.ETChannelThreshold, supervisor.Config.ETChannelGrowthFactor)
	supervisor.TLChannel = channel.NewManagedChannel("TLChannel", supervisor.Config.TLChannelThreshold, supervisor.Config.TLChannelGrowthFactor)

	if helper != nil {
		supervisor.helper = helper
	} else {
		// TODO : I could not be arsed, clean it up later, this stinks
		panic("helper can not be nil")
	}

	if metadata != nil {
		supervisor.Metadata = NewMetadata(metadata)
	} else {
		supervisor.Metadata = NewMetadata(nil)
	}
	return supervisor
}

func NewCustomSupervisor(clusterImplementation cluster.Cluster, config *cluster.Config, metadata map[string]string, helper cluster.H) *Supervisor {
	supervisor := new(Supervisor)

	/**
	 * Note: we may wish to dynamically modify the threshold and growth-factor rates
	 *       used by the managed channels to vary how provisioning of new transform and
	 *       load goroutines are created. This allows us to create an autonomous system
	 *       that "self improves" if the output of the monitor is looped back
	 */

	supervisor.State = UnTouched
	supervisor.group = clusterImplementation
	supervisor.Config = *config // copy config
	supervisor.Stats = cluster.NewStatistics()
	supervisor.ETChannel = channel.NewManagedChannel("ETChannel", config.ETChannelThreshold, config.ETChannelGrowthFactor)
	supervisor.TLChannel = channel.NewManagedChannel("TLChannel", config.TLChannelThreshold, config.TLChannelGrowthFactor)

	if helper != nil {
		supervisor.helper = helper
	} else {
		// TODO : fix
		panic("helper cannot be nil")
	}

	if metadata != nil {
		supervisor.Metadata = NewMetadata(metadata)
	} else {
		supervisor.Metadata = NewMetadata(nil)
	}

	return supervisor
}
