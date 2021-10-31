package app

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoAppProvider struct {
	BaseFolder string
}

func (y *YogoAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewYogoApp
}

func (y *YogoAppProvider) Boot(container framework.Container) error {
	return nil
}

func (y *YogoAppProvider) IsDefer() bool {
	return false
}

func (y *YogoAppProvider) Name() string {
	return contract.AppKey
}

func (y *YogoAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, y.BaseFolder}
}
