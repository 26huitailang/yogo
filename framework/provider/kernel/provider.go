package kernel

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
)

type YogoKernelProvider struct {
	HttpEngine *gin.Engine
}

func (y *YogoKernelProvider) Register(container framework.Container) framework.NewInstance {
	return NewYogoKernelService
}

func (y *YogoKernelProvider) Boot(container framework.Container) error {
	if y.HttpEngine == nil {
		y.HttpEngine = gin.Default()
	}
	y.HttpEngine.SetContainer(container)
	return nil
}

func (y *YogoKernelProvider) IsDefer() bool {
	return false
}

func (y *YogoKernelProvider) Name() string {
	return contract.KernelKey
}

func (y *YogoKernelProvider) Params(container framework.Container) []interface{} {
	return []interface{}{y.HttpEngine}
}
