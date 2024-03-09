package common

import "github.com/GabeCordo/processor-framework/processor/components/cluster"

type InterruptEvent uint8

const (
	Shutdown InterruptEvent = 0
	Panic                   = 1
)

type ProvisionerConfig struct {
	Debug string
}

type ProvisionerAction uint8

const (
	ProvisionerModuleGet ProvisionerAction = iota
	ProvisionerSupervisorGet
	ProvisionerSupervisorCreate
)

type ProvisionerSource string

const (
	Core ProvisionerSource = "core"
	User                   = "user"
)

type ProvisionerRequest struct {
	Action   ProvisionerAction
	Source   ProvisionerSource
	Module   string
	Cluster  string
	Config   *cluster.Config
	Metadata map[string]string
	Path     string
	Nonce    uint32
}

type ProvisionerResponse struct {
	Success bool
	Error   error
	Data    any
	Nonce   uint32
}
