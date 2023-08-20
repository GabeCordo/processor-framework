package controllers

import (
	"fmt"
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/mango-go/processor/threads"
	"github.com/GabeCordo/mango-go/processor/threads/common"
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
		fmt.Printf("[✓] the core has been initialized (%s)\n", threads.DefaultModulesFolder)
	}

	if _, err := os.Stat(threads.DefaultProcessorFolder); err != nil {
		fmt.Printf("[x] the processor folder is missing (%s)\n", threads.DefaultProcessorFolder)
	} else {
		fmt.Printf("[✓] the processor folder exists (%s)\n", threads.DefaultProcessorFolder)
	}

	if _, err := os.Stat(threads.DefaultModulesFolder); err != nil {
		fmt.Printf("[x] the modules folder is missing (%s)\n", threads.DefaultModulesFolder)
	} else {
		fmt.Printf("[✓] the modules folder exists (%s)\n", threads.DefaultModulesFolder)
	}

	// TODO : verify processor config file
	if _, err := os.Stat(threads.DefaultProcessorConfig); err != nil {
		fmt.Printf("[x] the processor config is missing (%s)\n", threads.DefaultProcessorFolder)
		return commandline.Terminate
	}

	src, err := os.Open(threads.DefaultProcessorConfig)
	if err != nil {
		fmt.Printf("[x] failed to read the processor config (%s)\n", threads.DefaultProcessorConfig)
		return commandline.Terminate
	}

	cfg := &common.Config{}
	err = yaml.NewDecoder(src).Decode(cfg)
	if err != nil {
		fmt.Printf("[x] the processor config is corrupt (%s)\n", threads.DefaultProcessorConfig)
	} else {
		fmt.Printf("[✓] the processor config is valid (%s)\n", threads.DefaultProcessorConfig)
	}

	return commandline.Terminate
}
