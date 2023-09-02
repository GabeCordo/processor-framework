package provisioner

import (
	"errors"
)

func NewProvisioner() *Provisioner {
	provisioner := new(Provisioner)

	provisioner.modules = make(map[string]*ModuleWrapper)

	defaultFrameworkModule := NewModuleWrapper()
	defaultFrameworkModule.Identifier = DefaultFrameworkModule
	defaultFrameworkModule.Version = 1.0
	defaultFrameworkModule.Mount()
	provisioner.modules[DefaultFrameworkModule] = defaultFrameworkModule

	return provisioner
}

func (provisioner *Provisioner) ModuleExists(moduleName string) bool {
	provisioner.mutex.RLock()
	defer provisioner.mutex.RUnlock()

	_, found := provisioner.modules[moduleName]
	return found
}

func (provisioner *Provisioner) GetModules() []*ModuleWrapper {

	provisioner.mutex.RLock()
	defer provisioner.mutex.RUnlock()

	modules := make([]*ModuleWrapper, 0)
	for _, moduleWrapper := range provisioner.modules {
		modules = append(modules, moduleWrapper)
	}

	return modules
}

func (provisioner *Provisioner) GetModule(moduleName string) (instance *ModuleWrapper, found bool) {
	provisioner.mutex.RLock()
	defer provisioner.mutex.RUnlock()

	instance, found = provisioner.modules[moduleName]
	if !found {
		return nil, false
	}

	if instance.MarkForDeletion {
		return nil, false
	}

	return instance, found
}

func (provisioner *Provisioner) AddModule(identifier string) error {

	provisioner.mutex.Lock()
	defer provisioner.mutex.Unlock()

	if _, found := provisioner.modules[identifier]; found {
		return errors.New("module already exists")
	}

	mod := NewModuleWrapper()
	mod.Mounted = true
	mod.Identifier = identifier
	mod.Version = 1.0
	provisioner.modules[identifier] = mod

	return nil
}

func (provisioner *Provisioner) DeleteModule(identifier string) (deleted, markedForDeletion, found bool) {

	provisioner.mutex.Lock()
	defer provisioner.mutex.Unlock()

	deleted = false

	if moduleWrapper, foundModule := provisioner.modules[identifier]; foundModule {
		found = true

		provisioner.modules[identifier].MarkForDeletion = true
		markedForDeletion = true

		if moduleWrapper.CanDelete() {
			delete(provisioner.modules, identifier)
			deleted = true
		} else {
		}
	} else {
		found = false
	}

	return deleted, markedForDeletion, found
}
