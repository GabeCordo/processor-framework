package controllers

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/keitt/processor"
	"gopkg.in/yaml.v3"
	"os"
)

type StartCommand struct {
}

func (command StartCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	// Note: the processor is being launched in standalone mode
	//		-> no connection is being made to a core
	//		-> all stream processes are started on launch
	//		-> all batch processes can be run via repl

	f, err := os.Open(DefaultProcessorConfig)
	if err != nil {
		panic(err)
	}

	cfg := &processor.Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		panic(err)
	}

	processor, err := processor.New(cfg)
	if err != nil {
		panic(err)
	}

	processor.Run()

	return commandline.Terminate
}
