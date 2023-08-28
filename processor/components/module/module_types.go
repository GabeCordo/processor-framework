package module

import (
	"github.com/GabeCordo/mango/core/interfaces/module"
	"plugin"
)

type Module struct {
	Plugin *plugin.Plugin
	Config *module.Config
}

type RemoteModule struct {
	Path string
}
