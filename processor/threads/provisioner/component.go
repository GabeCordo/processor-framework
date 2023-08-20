package provisioner

import "github.com/GabeCordo/mango-go/processor/components/provisioner"

var instance *provisioner.Provisioner

func GetProvisionerInstance() *provisioner.Provisioner {

	if instance == nil {
		instance = provisioner.NewProvisioner()
	}
	return instance
}
