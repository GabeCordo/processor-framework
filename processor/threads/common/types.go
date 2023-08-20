package common

import "github.com/GabeCordo/mango/components/cluster"

type ProvisionerAction string

const (
	ProvisionerCreateModule     ProvisionerAction = "create"
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
	Nonce   uint32
}
