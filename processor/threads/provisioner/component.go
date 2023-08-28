package provisioner

import "github.com/GabeCordo/keitt/processor/components/provisioner"

var instance *provisioner.Provisioner

func GetProvisionerInstance() *provisioner.Provisioner {

	if instance == nil {
		instance = provisioner.NewProvisioner()
	}
	return instance
}
