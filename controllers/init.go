package controllers

import (
	"fmt"
	"github.com/GabeCordo/commandline"
	"github.com/GabeCordo/keitt/processor"
	"github.com/GabeCordo/mango/core/threads/common"
	"gopkg.in/yaml.v3"
	"os"
)

type InitCommand struct {
}

func (cmd InitCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	if _, err := os.Stat(common.DefaultFrameworkFolder); err == nil {
		fmt.Printf("[-] the framework folder has already been initialized (%s)\n", common.DefaultFrameworkFolder)
	} else {
		fmt.Println(os.IsExist(err))
		if err = os.Mkdir(common.DefaultFrameworkFolder, 0700); err != nil {
			fmt.Printf("[x] the root processor folder could not be initialized %s (%s)\n", err.Error(), common.DefaultFrameworkFolder)
			return commandline.Terminate
		} else {
			fmt.Printf("[✓] the root processor folder has been initialized (%s)\n", common.DefaultFrameworkFolder)
		}
	}

	if _, err := os.Stat(DefaultProcessorFolder); err == nil {
		fmt.Printf("[-] the processor folder has already been initialized (%s)\n", DefaultProcessorFolder)
	} else {
		if err := os.Mkdir(DefaultProcessorFolder, 0700); err != nil {
			fmt.Printf("[-] failed to create the processor folder (%s)\n", DefaultProcessorFolder)
			return commandline.Terminate
		} else {
			fmt.Printf("[✓] the processor folder has been initialized (%s)\n", DefaultProcessorFolder)
		}
	}

	if _, err := os.Stat(DefaultProcessorConfig); err == nil {
		fmt.Printf("[-] the global config has already been initialized (%s)\n", DefaultProcessorConfig)
	} else {
		dst, err := os.Create(DefaultProcessorConfig)
		if err != nil {
			fmt.Printf("[x] failed to create the processor config %s\n", DefaultProcessorConfig)
			return commandline.Terminate
		}

		defaultConfig := processor.Config{
			Name:    "a",
			Debug:   true,
			Timeout: 2,
			Net: struct {
				Host string `yaml:"host"`
				Port int    `yaml:"port"`
			}(struct {
				Host string
				Port int
			}{Host: "localhost", Port: 5023}),
		}

		b, err := yaml.Marshal(defaultConfig)
		if err != nil {
			fmt.Printf("[x] failed to generate the processor config %s\n", DefaultProcessorConfig)
			return commandline.Terminate
		}

		if _, err := dst.Write(b); err != nil {
			fmt.Printf("[x] failed to create the processor config %s\n", DefaultProcessorConfig)
		} else {
			fmt.Printf("[✓] created the processor config %s\n", DefaultProcessorConfig)
		}
	}

	return commandline.Terminate
}
