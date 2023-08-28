package helper

import (
	"github.com/GabeCordo/mango/api"
)

type Helper interface {
	IsDebugEnabled() bool
	SaveToCache(data any) (string, error)
	LoadFromCache(identifier string) (any, error)
	Log(message string)
	Warning(message string)
	Fatal(message string)
}

type UsesHelper interface {
	SetHelper(helper Helper)
}

type StandardHelper struct {
	host    string
	module  string
	cluster string
}

func NewStandardHelper(module, cluster string) (*StandardHelper, error) {
	helper := new(StandardHelper)

	// TODO : setup core host link

	helper.module = module
	helper.cluster = cluster

	return helper, nil
}

func (helper StandardHelper) IsDebugEnabled() bool {
	// TODO : add debug check
	return api.IsDebugEnabled(helper.host)
}

func (helper StandardHelper) SaveToCache(data any) (string, error) {

	return api.Cache(helper.host, data)
}

func (helper StandardHelper) LoadFromCache(identifier string) (any, error) {

	return api.GetFromCache(helper.host, identifier)
}

func (helper StandardHelper) Log(message string) {

	api.Log(helper.host, message)
}

func (helper StandardHelper) Warning(message string) {

	api.LogWarn(helper.host, message)
}

func (helper StandardHelper) Fatal(message string) {

	api.LogError(helper.host, message)
}
