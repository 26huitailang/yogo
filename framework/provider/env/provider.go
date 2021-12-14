package env

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type YogoEnvProvider struct {
	Folder string
}

func (provider *YogoEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewYogoEnv
}

func (provider *YogoEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

func (provider *YogoEnvProvider) IsDefer() bool {
	return false
}

func (provider *YogoEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

func (provider *YogoEnvProvider) Name() string {
	return contract.EnvKey
}
