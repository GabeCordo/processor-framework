package provisioner

import (
	"errors"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
	"sync"
)

type Config struct {
	Debug      bool
	Timeout    float64
	Standalone bool
	Core       string
}

type Thread struct {
	Config *Config

	Interrupt chan<- threads.InterruptEvent // Upon completion or failure an interrupt can be raised

	C1 chan common.ProvisionerRequest    // Supervisor is receiving threads from the http_thread
	C2 chan<- common.ProvisionerResponse // Supervisor is sending responses to the http_thread
	
	logger *utils.Logger

	accepting   bool
	listenersWg sync.WaitGroup
	requestWg   sync.WaitGroup
}

func NewThread(cfg *Config, logger *utils.Logger, channels ...interface{}) (*Thread, error) {
	provisioner := new(Thread)
	var ok bool

	provisioner.Interrupt, ok = (channels[0]).(chan threads.InterruptEvent)
	if !ok {
		return nil, errors.New("expected type 'chan InterruptEvent' in index 0")
	}
	provisioner.C1, ok = (channels[1]).(chan common.ProvisionerRequest)
	if !ok {
		return nil, errors.New("expected type 'chan ProvisionerRequest' in index 1")
	}
	provisioner.C2, ok = (channels[2]).(chan common.ProvisionerResponse)
	if !ok {
		return nil, errors.New("expected type 'chan ProvisionerResponse' in index 2")
	}

	if logger == nil {
		return nil, errors.New("expected non nil *utils.Logger type")
	}
	provisioner.logger = logger

	if cfg == nil {
		return nil, errors.New("expected no nil *provisioner.Config type")
	}
	provisioner.Config = cfg

	provisioner.logger.SetColour(utils.Orange)

	return provisioner, nil
}
