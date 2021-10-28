package app

import (
	"errors"
	"github.com/26huitailang/yogo/framework"
)

type YogoAppProvider struct {
	BaseFolder string
}

func (y *YogoAppProvider) Register(container framework.Container) framework.NewInstance {
	panic("implement me")
}

func (y *YogoAppProvider) Boot(container framework.Container) error {
	panic("implement me")
}

func (y *YogoAppProvider) IsDefer() bool {
	panic("implement me")
}

func (y *YogoAppProvider) Name() string {
	panic("implement me")
}

func (y *YogoAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, y.BaseFolder}
}

func NewYogoApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return &YogoApp{baseFolder: baseFolder, container: container}, nil
}
