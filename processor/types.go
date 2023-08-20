package processor

import (
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango-go/processor/threads/http"
	"github.com/GabeCordo/mango-go/processor/threads/provisioner"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
)

type Config struct {
	Name               string  `yaml:"name"`
	Debug              bool    `yaml:"debug"`
	MaxWaitForResponse float64 `yaml:"max-wait-for-response"`
	Core               string  `yaml:"core"`
	Net                struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"net"`
	StandaloneMode bool `yaml:"standalone-mode"`
}

type Processor struct {
	HttpThread  *http.Thread
	Provisioner *provisioner.Thread

	Interrupt chan threads.InterruptEvent
	C1        chan common.ProvisionerRequest
	C2        chan common.ProvisionerResponse

	Config *Config

	Logger *utils.Logger
}
