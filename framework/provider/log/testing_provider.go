package log

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoTestingLogProvider struct {
}

// Register registe a new function for make a service instance
func (provider *YogoTestingLogProvider) Register(c framework.Container) framework.NewInstance {
	return NewYogoTestingLog
}

// Boot will called when the service instantiate
func (provider *YogoTestingLogProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YogoTestingLogProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YogoTestingLogProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

// / Name define the name for this service
func (provider *YogoTestingLogProvider) Name() string {
	return contract.LogKey
}
