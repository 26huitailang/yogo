package distributed

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

type LocalDistributedProvider struct{}

func (l *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

func (l *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

func (l *LocalDistributedProvider) IsDefer() bool {
	return false
}

func (l *LocalDistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (l *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
