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

type InitCommand struct {
}

func (cmd InitCommand) Run(cli *commandline.CommandLine) commandline.TerminateOnCompletion {

	if _, err := os.Stat(core.DefaultFrameworkFolder); err == nil {
		fmt.Printf("[-] the framework folder has already been initialized (%s)\n", core.DefaultFrameworkFolder)
	} else {
		fmt.Println(os.IsExist(err))
		if err = os.Mkdir(core.DefaultFrameworkFolder, 0700); err != nil {
			fmt.Printf("[x] the root processor folder could not be initialized %s (%s)\n", err.Error(), core.DefaultFrameworkFolder)
			return commandline.Terminate
		} else {
			fmt.Printf("[✓] the root processor folder has been initialized (%s)\n", core.DefaultFrameworkFolder)
		}
	}

	if _, err := os.Stat(threads.DefaultProcessorFolder); err == nil {
		fmt.Printf("[-] the processor folder has already been initialized (%s)\n", core.DefaultFrameworkFolder)
	} else {
		if err := os.Mkdir(threads.DefaultProcessorFolder, 0700); err != nil {
			fmt.Printf("[-] failed to create the processor folder (%s)\n", core.DefaultFrameworkFolder)
			return commandline.Terminate
		} else {
			fmt.Printf("[✓] the processor folder has been initialized (%s)\n", threads.DefaultModulesFolder)
		}
	}

	if _, err := os.Stat(threads.DefaultModulesFolder); err == nil {
		fmt.Printf("[-] the modules folder has already been initialized (%s)\n", core.DefaultFrameworkFolder)
	} else {
		if err := os.Mkdir(threads.DefaultModulesFolder, 0700); err != nil {
			fmt.Printf("[x] failed to create %s directory %s\n", threads.DefaultModulesFolder, err.Error())
			return commandline.Terminate
		} else {
			fmt.Printf("[✓] created modules folder %s\n", threads.DefaultModulesFolder)
		}
	}

	if _, err := os.Stat(threads.DefaultProcessorConfig); err == nil {
		fmt.Printf("[-] the global config has already been initialized (%s)\n", core.DefaultFrameworkFolder)
	} else {
		dst, err := os.Create(threads.DefaultProcessorConfig)
		if err != nil {
			fmt.Printf("[x] failed to create the processor config %s\n", threads.DefaultProcessorConfig)
			return commandline.Terminate
		}

		defaultConfig := common.Config{
			Name:               "a",
			Debug:              true,
			MaxWaitForResponse: 2,
			Net: struct {
				Host string `yaml:"host"`
				Port int    `yaml:"port"`
			}(struct {
				Host string
				Port int
			}{Host: "localhost", Port: 5023}),
			Path: threads.DefaultProcessorConfig,
		}

		b, err := yaml.Marshal(defaultConfig)
		if err != nil {
			fmt.Printf("[x] failed to generate the processor config %s\n", threads.DefaultProcessorConfig)
			return commandline.Terminate
		}

		if _, err := dst.Write(b); err != nil {
			fmt.Printf("[x] failed to create the processor config %s\n", threads.DefaultProcessorConfig)
		} else {
			fmt.Printf("[✓] created the processor config %s\n", threads.DefaultProcessorConfig)
		}
	}

	return commandline.Terminate
}
