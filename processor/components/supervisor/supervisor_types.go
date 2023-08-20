package supervisor

import (
	"github.com/GabeCordo/mango-go/processor/components/channel"
	"github.com/GabeCordo/mango/components/cluster"
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
	Unknown             = "unknown"
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
	ETChannel *channel.ManagedChannel
	TLChannel *channel.ManagedChannel

	loadWaitGroup sync.WaitGroup
	waitGroup     sync.WaitGroup
	mutex         sync.RWMutex
}
