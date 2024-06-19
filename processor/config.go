package processor

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func NewConfig(name string) *Config {
	config := new(Config)
	config.Name = name
	config.Net.External.Host = "localhost"
	config.Net.External.Port = 5023
	config.Net.Internal.Host = "localhost"
	config.Net.Internal.Port = 5023
	config.Debug = true
	return config
}

func ConfigFromYAML(config *Config, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// file does not exist
		log.Println("foo")
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
