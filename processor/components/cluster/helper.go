package cluster

import (
	"github.com/GabeCordo/mango/api"
)

type StandardHelper struct {
	host       string
	module     string
	cluster    string
	supervisor uint64
}

func NewHelper(host, module, cluster string, supervisor uint64) *StandardHelper {
	helper := new(StandardHelper)

	// TODO : setup core host link

	helper.host = host
	helper.module = module
	helper.cluster = cluster
	helper.supervisor = supervisor

	return helper
}

func (helper StandardHelper) IsDebugEnabled() bool {
	// TODO : add debug check
	return api.IsDebugEnabled(helper.host)
}

func (helper StandardHelper) SaveToCache(data string) (string, error) {

	return api.Cache(helper.host, "", data)
}

func (helper StandardHelper) LoadFromCache(identifier string) (string, error) {

	return api.GetFromCache(helper.host, identifier)
}

func (helper StandardHelper) Log(message string) error {

	return api.Log(helper.host, helper.supervisor, message)
}

func (helper StandardHelper) Warning(message string) error {

	return api.LogWarn(helper.host, helper.supervisor, message)
}

func (helper StandardHelper) Fatal(message string) error {

	return api.LogError(helper.host, helper.supervisor, message)
}
