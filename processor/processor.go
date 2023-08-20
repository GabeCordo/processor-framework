package processor

import (
	"errors"
	"fmt"
	provisioner2 "github.com/GabeCordo/mango-go/processor/components/provisioner"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango-go/processor/threads/http"
	"github.com/GabeCordo/mango-go/processor/threads/provisioner"
	"github.com/GabeCordo/mango/proxy"
	"github.com/GabeCordo/mango/threads"
	"github.com/GabeCordo/mango/utils"
	"os"
	"os/signal"
	"syscall"
)

func New(cfg *Config) (*Processor, error) {
	processor := new(Processor)

	if cfg == nil {
		panic(errors.New("the config passed to processor.New cannot be nil"))
	}
	processor.Config = cfg

	processor.Interrupt = make(chan threads.InterruptEvent, 1)
	processor.C1 = make(chan common.ProvisionerRequest, 10)
	processor.C2 = make(chan common.ProvisionerResponse, 10)

	httpConfig := &http.Config{
		Debug:   cfg.Debug,
		Timeout: cfg.MaxWaitForResponse,
		Net:     fmt.Sprintf("%s:%d", cfg.Net.Host, cfg.Net.Port),
	}
	httpLogger, err := utils.NewLogger(utils.HttpProcessor, &cfg.Debug)
	if err != nil {
		return nil, err
	}
	processor.HttpThread, err = http.NewThread(httpConfig, httpLogger,
		processor.Interrupt, processor.C1, processor.C2)

	provisionerConfig := &provisioner.Config{
		Debug:      true,
		Timeout:    cfg.MaxWaitForResponse,
		Standalone: cfg.StandaloneMode,
		Core:       cfg.Core,
	}
	provisionerLogger, err := utils.NewLogger(utils.Provisioner, &processor.Config.Debug)
	if err != nil {
		return nil, err
	}
	processor.Provisioner, err = provisioner.NewThread(provisionerConfig, provisionerLogger,
		processor.Interrupt, processor.C1, processor.C2)
	if err != nil {
		return nil, err
	}

	processorLogger, err := utils.NewLogger(utils.Undefined, &processor.Config.Debug)
	if err != nil {
		return nil, err
	}
	processor.Logger = processorLogger

	return processor, nil
}

func (processor *Processor) Run() {

	if !processor.Config.StandaloneMode {
		err := proxy.ConnectToCore(processor.Config.Core)
		if err != nil {
			panic(err)
		}
	}

	processor.Logger.SetColour(utils.Purple)

	if processor.Config.Debug {
		if processor.Config.StandaloneMode {
			processor.Logger.Println("running in STANDALONE mode")
		} else {
			processor.Logger.Println("running in CONNECTED mode")
		}
	}

	processor.Provisioner.Setup()
	if processor.Config.Debug {
		processor.Logger.Println("started provisioner thread")
	}
	go processor.Provisioner.Start()

	processor.HttpThread.Setup()
	if processor.Config.Debug {
		processor.Logger.Println("started http processor thread")
	}
	go processor.HttpThread.Start()

	go processor.repl()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		processor.Interrupt <- threads.Panic
	case interrupt := <-processor.Interrupt:
		switch interrupt {
		case threads.Panic:
			processor.Logger.Printf("[IO] %s\n", " encountered panic")
		default: // shutdown
			processor.Logger.Printf("[IO] %s\n", " shutting down")
		}
	}

	processor.Logger.SetColour(utils.Red)

	processor.HttpThread.Teardown()
	if processor.Config.Debug {
		processor.Logger.Println("http processor thread shutdown")
	}

	processor.Provisioner.Teardown()
	if processor.Config.Debug {
		processor.Logger.Println("provisioner thread shutdown")
	}
}

func (processor *Processor) Module(name string) *provisioner2.ModuleWrapper {

	if _, found := provisioner.GetProvisionerInstance().GetModule(name); !found {
		provisioner.GetProvisionerInstance().AddModule(name)
	}

	mod, _ := provisioner.GetProvisionerInstance().GetModule(name)
	return mod
}
