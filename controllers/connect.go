package controllers

import (
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/keitt/processor"
	"gopkg.in/yaml.v3"
	"os"
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

	f, err := os.Open(DefaultProcessorConfig)
	if err != nil {
		panic(err)
	}

	cfg := &processor.Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		panic(err)
	}

	cfg.Core = coreEndpoint
	cfg.StandaloneMode = false

	processor, err := processor.New(cfg)
	if err != nil {
		panic(err)
	}

	processor.Run()

	return commandline.Terminate
}
