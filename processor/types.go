package processor

import (
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/keitt/processor/threads/http"
	"github.com/GabeCordo/keitt/processor/threads/provisioner"
	"github.com/GabeCordo/toolchain/logging"
)

type Module uint8

const (
	HttpProcessor Module = iota
	Provisioner
	Undefined
)

func (module Module) ToString() string {
	switch module {
	case HttpProcessor:
		return "http-processor"
	case Provisioner:
		return "provisioner"
	default:
		return "unknown"
	}
}

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

	Interrupt chan common.InterruptEvent
	C1        chan common.ProvisionerRequest
	C2        chan common.ProvisionerResponse

	Config *Config

	Logger *logging.Logger
}
