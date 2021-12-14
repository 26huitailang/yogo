package config

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoConfigProvider struct {
	c      framework.Container
	folder string
	env    string

	envMaps map[string]string
}

// Register registe a new function for make a service instance
func (provider *YogoConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewYogoConfig
}

// Boot will called when the service instantiate
func (provider *YogoConfigProvider) Boot(c framework.Container) error {
	provider.folder = c.MustMake(contract.AppKey).(contract.App).ConfigFolder()
	provider.envMaps = c.MustMake(contract.EnvKey).(contract.Env).All()
	provider.env = c.MustMake(contract.EnvKey).(contract.Env).AppEnv()
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YogoConfigProvider) IsDefer() bool {
	return true
}

// Params define the necessary params for NewInstance
func (provider *YogoConfigProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.folder, provider.envMaps, provider.env, provider.c}
}

/// Name define the name for this service
func (provider *YogoConfigProvider) Name() string {
	return contract.ConfigKey
}
