package controllers

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/mango-go/processor/threads"
	"github.com/GabeCordo/mango-go/processor/threads/common"
)

type StartCommand struct {
}

func (command StartCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	// Note: the processor is being launched in standalone mode
	//		-> no connection is being made to a core
	//		-> all stream processes are started on launch
	//		-> all batch processes can be run via repl

	config := common.GetConfigInstance(threads.DefaultProcessorConfig)
	if config == nil {
		panic("could not find processor config")
	}
	config.StandaloneMode = true

	processor, err := threads.NewProcessor(config)
	if err != nil {
		panic(err)
	}

	processor.Run()

	return commandline.Terminate
}
