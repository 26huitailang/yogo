package env

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoTestingEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *YogoTestingEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewYogoTestingEnv
}

// Boot will called when the service instantiate
func (provider *YogoTestingEnvProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YogoTestingEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YogoTestingEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

// / Name define the name for this service
func (provider *YogoTestingEnvProvider) Name() string {
	return contract.EnvKey
}
