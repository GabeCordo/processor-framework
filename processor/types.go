package processor

import (
	"errors"
	"fmt"
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/keitt/processor/threads/http"
	"github.com/GabeCordo/keitt/processor/threads/provisioner"
	processor_i "github.com/GabeCordo/mango/core/interfaces/processor"
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
		return "-"
	}
}

type Config struct {
	Name           string  `yaml:"name"`
	Debug          bool    `yaml:"debug"`
	StandaloneMode bool    `yaml:"standalone"`
	ReplMode       bool    `yaml:"repl"`
	Timeout        float64 `yaml:"timeout"`
	Core           string  `yaml:"core"`
	Net            struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"net"`
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

func New(cfg *Config) (*Processor, error) {
	processor := new(Processor)

	if cfg == nil {
		panic(errors.New("the config passed to processor.New cannot be nil"))
	}
	processor.Config = cfg

	processor.Interrupt = make(chan common.InterruptEvent, 1)
	processor.C1 = make(chan common.ProvisionerRequest, 10)
	processor.C2 = make(chan common.ProvisionerResponse, 10)

	httpConfig := &http.Config{
		Debug:   cfg.Debug,
		Timeout: cfg.Timeout,
		Net:     fmt.Sprintf("%s:%d", cfg.Net.Host, cfg.Net.Port),
	}
	httpLogger, err := logging.NewLogger(HttpProcessor.ToString(), &cfg.Debug)
	if err != nil {
		return nil, err
	}
	processor.HttpThread, err = http.NewThread(httpConfig, httpLogger,
		processor.Interrupt, processor.C1, processor.C2)

	provisionerConfig := &provisioner.Config{
		Debug:      true,
		Timeout:    cfg.Timeout,
		Standalone: cfg.StandaloneMode,
		Core:       cfg.Core,
		Processor:  processor_i.Config{Host: cfg.Net.Host, Port: cfg.Net.Port},
	}
	provisionerLogger, err := logging.NewLogger(Provisioner.ToString(), &processor.Config.Debug)
	if err != nil {
		return nil, err
	}
	processor.Provisioner, err = provisioner.NewThread(provisionerConfig, provisionerLogger,
		processor.Interrupt, processor.C1, processor.C2)
	if err != nil {
		return nil, err
	}

	processorLogger, err := logging.NewLogger(Undefined.ToString(), &processor.Config.Debug)
	if err != nil {
		return nil, err
	}
	processor.Logger = processorLogger

	return processor, nil
}
