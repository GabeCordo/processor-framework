package common

import "github.com/GabeCordo/keitt/processor/components/cluster"

type InterruptEvent uint8

const (
	Shutdown InterruptEvent = 0
	Panic                   = 1
)

type ProvisionerConfig struct {
	Debug string
}

type ProvisionerAction string

const (
	ProvisionerGetModules       ProvisionerAction = "get"
	ProvisionerCreateModule                       = "create"
	ProvisionerDeleteModule                       = "delete"
	ProvisionerCreateSupervisor                   = "provision"
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
