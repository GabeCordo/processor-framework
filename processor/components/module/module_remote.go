package module

import (
	"github.com/GabeCordo/mango/core/interfaces/module"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

func NewRemoteModule(path string) (*RemoteModule, error) {

	remoteModule := new(RemoteModule)
	remoteModule.Path = path

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return remoteModule, nil
}

func (remoteModule RemoteModule) Get() (*Module, error) {

	mod := new(Module)

	filepath.Walk(remoteModule.Path, func(path string, info os.FileInfo, err error) error {

		if info.Name() == "module.etl.yaml" {
			f, err := os.Open(path)
			if err != nil {
				log.Println(err)
				return err
			}

			cfg := new(module.Config)
			if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
				log.Println(err)
				return err
			} else {
				mod.Config = cfg
			}
		} else if filepath.Ext(info.Name()) == ".so" {
			if mod.Plugin, err = plugin.Open(path); err != nil {
				log.Println(err)
				return err
			}
		}

		return nil
	})

	if (mod.Plugin == nil) || (mod.Config == nil) {
		return nil, os.ErrExist
	}

	return mod, nil
}
