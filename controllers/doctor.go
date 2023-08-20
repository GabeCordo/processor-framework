package controllers

import (
	"fmt"
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/mango-go/processor"
	"github.com/GabeCordo/mango/core"
	"gopkg.in/yaml.v3"
	"os"
)

type DoctorCommand struct {
}

func (cmd DoctorCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	// TODO : verify the core folder is setup
	if _, err := os.Stat(core.DefaultFrameworkFolder); err != nil {
		fmt.Printf("[x] the core has not been initialized (%s)\n", core.DefaultFrameworkFolder)
	} else {
		fmt.Printf("[✓] the core has been initialized (%s)\n", core.DefaultFrameworkFolder)
	}

	if _, err := os.Stat(DefaultProcessorFolder); err != nil {
		fmt.Printf("[x] the processor folder is missing (%s)\n", DefaultProcessorFolder)
	} else {
		fmt.Printf("[✓] the processor folder exists (%s)\n", DefaultProcessorFolder)
	}

	if _, err := os.Stat(DefaultModulesFolder); err != nil {
		fmt.Printf("[x] the modules folder is missing (%s)\n", DefaultModulesFolder)
	} else {
		fmt.Printf("[✓] the modules folder exists (%s)\n", DefaultModulesFolder)
	}

	// TODO : verify processor config file
	if _, err := os.Stat(DefaultProcessorConfig); err != nil {
		fmt.Printf("[x] the processor config is missing (%s)\n", DefaultProcessorFolder)
		return commandline.Terminate
	}

	src, err := os.Open(DefaultProcessorConfig)
	if err != nil {
		fmt.Printf("[x] failed to read the processor config (%s)\n", DefaultProcessorConfig)
		return commandline.Terminate
	}

	cfg := &processor.Config{}
	err = yaml.NewDecoder(src).Decode(cfg)
	if err != nil {
		fmt.Printf("[x] the processor config is corrupt (%s)\n", DefaultProcessorConfig)
	} else {
		fmt.Printf("[✓] the processor config is valid (%s)\n", DefaultProcessorConfig)
	}

	return commandline.Terminate
}
