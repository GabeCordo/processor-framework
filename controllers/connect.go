package controllers

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/mango-go/processor/threads"
	"github.com/GabeCordo/mango-go/processor/threads/common"
)

type ConnectCommand struct {
}

func (cmd ConnectCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	// Note: the processor is being run in connected mode
	//		-> the processor is establishing a "session" with a core server
	//		-> streams are provisioned by the core upon registering the module
	//		-> batch are provisioned by the operator through the core

	// TODO : verify establish a connection with the core
	coreEndpoint := cli.NextArg()
	if coreEndpoint == commandline.FinalArg {
		panic("missing core endpoint host")
	}

	config := common.GetConfigInstance(threads.DefaultProcessorConfig)
	if config == nil {
		panic("could not find processor config")
	}
	config.Core = coreEndpoint
	config.StandaloneMode = false

	processor, err := threads.NewProcessor(config)
	if err != nil {
		panic(err)
	}

	processor.Run()

	return commandline.Terminate
}
