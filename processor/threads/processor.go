package threads

import (
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

func NewProcessor(config *common.Config) (*Processor, error) {
	processor := new(Processor)

	processor.Interrupt = make(chan threads.InterruptEvent, 1)
	processor.C1 = make(chan common.ProvisionerRequest, 10)
	processor.C2 = make(chan common.ProvisionerResponse, 10)

	httpLogger, err := utils.NewLogger(utils.HttpProcessor, &common.GetConfigInstance().Debug)
	if err != nil {
		return nil, err
	}
	processor.HttpThread, err = http.NewThread(httpLogger,
		processor.Interrupt, processor.C1, processor.C2)

	provisionerLogger, err := utils.NewLogger(utils.Provisioner, &common.GetConfigInstance().Debug)
	if err != nil {
		return nil, err
	}
	processor.Provisioner, err = provisioner.NewThread(provisionerLogger, DefaultModulesFolder,
		processor.Interrupt, processor.C1, processor.C2)
	if err != nil {
		return nil, err
	}

	processorLogger, err := utils.NewLogger(utils.Undefined, &common.GetConfigInstance().Debug)
	if err != nil {
		return nil, err
	}
	processor.Logger = processorLogger

	return processor, nil
}

func (processor *Processor) Run() {

	if !common.GetConfigInstance().StandaloneMode {
		err := proxy.ConnectToCore(common.GetConfigInstance().Core)
		if err != nil {
			panic(err)
		}
	}

	processor.Logger.SetColour(utils.Purple)

	if common.GetConfigInstance().Debug {
		if common.GetConfigInstance().StandaloneMode {
			processor.Logger.Println("running in STANDALONE mode")
		} else {
			processor.Logger.Println("running in CONNECTED mode")
		}
	}

	processor.Provisioner.Setup()
	if common.GetConfigInstance().Debug {
		processor.Logger.Println("started provisioner thread")
	}
	go processor.Provisioner.Start()

	processor.HttpThread.Setup()
	if common.GetConfigInstance().Debug {
		processor.Logger.Println("started http processor thread")
	}
	go processor.HttpThread.Start()

	go repl()

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
	if common.GetConfigInstance().Debug {
		processor.Logger.Println("http processor thread shutdown")
	}

	processor.Provisioner.Teardown()
	if common.GetConfigInstance().Debug {
		processor.Logger.Println("provisioner thread shutdown")
	}
}

func (processor *Processor) Module(path string) error {

	return common.RegisterModule(processor.HttpThread.C1, processor.HttpThread.ProvisionerResponseTable, path)
}
