package trace

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *YogoTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewYogoTraceService
}

// Boot will called when the service instantiate
func (provider *YogoTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *YogoTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *YogoTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

// / Name define the name for this service
func (provider *YogoTraceProvider) Name() string {
	return contract.TraceKey
}
