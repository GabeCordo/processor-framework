package processor

import (
	provisioner2 "github.com/GabeCordo/keitt/processor/components/provisioner"
	"github.com/GabeCordo/keitt/processor/threads/common"
	"github.com/GabeCordo/keitt/processor/threads/provisioner"
	"github.com/GabeCordo/mango/api"
	processor_i "github.com/GabeCordo/mango/core/interfaces/processor"
	"github.com/GabeCordo/toolchain/logging"
	"os"
	"os/signal"
	"syscall"
)

func (processor *Processor) Run() {

	cfg := &processor_i.Config{Host: processor.Config.Net.Host, Port: processor.Config.Net.Port}

	if !processor.Config.StandaloneMode {
		err := api.ConnectToCore(processor.Config.Core, cfg)
		if err != nil {
			panic(err)
		}
	}

	processor.Logger.SetColour(logging.Purple)

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
		processor.Interrupt <- common.Panic
	case interrupt := <-processor.Interrupt:
		switch interrupt {
		case common.Panic:
			processor.Logger.Printf("[IO] %s\n", " encountered panic")
		default: // shutdown
			processor.Logger.Printf("[IO] %s\n", " shutting down")
		}
	}

	processor.Logger.SetColour(logging.Red)

	processor.HttpThread.Teardown()
	if processor.Config.Debug {
		processor.Logger.Println("http processor thread shutdown")
	}

	processor.Provisioner.Teardown()
	if processor.Config.Debug {
		processor.Logger.Println("provisioner thread shutdown")
	}

	if !processor.Config.StandaloneMode {
		err := api.DisconnectFromCore(processor.Config.Core, cfg)
		if err != nil {
			panic(err)
		}
	}
}

func (processor *Processor) Module(name string) *provisioner2.ModuleWrapper {

	if _, found := provisioner.GetProvisionerInstance().GetModule(name); !found {
		provisioner.GetProvisionerInstance().AddModule(name)
	}

	mod, _ := provisioner.GetProvisionerInstance().GetModule(name)
	return mod
}
