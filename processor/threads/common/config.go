package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Config struct {
	Name               string  `yaml:"name"`
	Debug              bool    `yaml:"debug"`
	MaxWaitForResponse float64 `yaml:"max-wait-for-response"`
	Core               string  `yaml:"core"`
	Net                struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"net"`
	Path           string `yaml:"path"`
	StandaloneMode bool   `yaml:"standalone-mode"`
}

func NewConfig(name string) *Config {
	config := new(Config)
	config.Name = name
	config.Net.Host = "localhost"
	config.Net.Port = 5023
	config.Debug = true
	return config
}

func yamlToETLConfig(config *Config, path string) error {
	if _, err := os.Stat(path); err != nil {
		// file does not exist
		log.Println(err)
		return err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		// error reading the file
		log.Println(err)
		return err
	}

	err = yaml.Unmarshal([]byte(file), config)
	if err != nil {
		// the file is not a JSON or is a malformed (fields missing) core
		log.Println(err)
		return err
	}

	return nil
}

var (
	configLock     = &sync.Mutex{}
	configInstance *Config
)

func GetConfigInstance(configPath ...string) *Config {
	configLock.Lock()
	defer configLock.Unlock()

	/* if this is the first time the common is being loaded the developer
	   needs to pass in a configPath to load the common instance from
	*/
	if (configInstance == nil) && (len(configPath) < 1) {
		return nil
	}

	if configInstance == nil {
		configInstance = NewConfig("test")

		if err := yamlToETLConfig(configInstance, configPath[0]); err == nil {
			// the configPath we found the common for future reference
			configInstance.Path = configPath[0]
			// if the MaxWaitForResponse is not set, then simply default to 2.0
			if configInstance.MaxWaitForResponse == 0 {
				configInstance.MaxWaitForResponse = 2
			}
		} else {
			log.Println("(!) the etl configuration file can either not be found or is corrupted")
			log.Fatal(fmt.Sprintf("%s was not a valid common path\n", configPath))
		}
	}

	return configInstance
}
