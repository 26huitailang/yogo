package demo

import (
	"fmt"
	"github.com/26huitailang/yogo/framework"
)

// DemoServiceProvider to provide demo service
type DemoServiceProvider struct{}

func (sp *DemoServiceProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

func (sp *DemoServiceProvider) Boot(c framework.Container) error {
	fmt.Println("demo service boot")
	return nil
}

func (sp *DemoServiceProvider) IsDefer() bool {
	return true
}

func (sp *DemoServiceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *DemoServiceProvider) Name() string {
	return DemoKey
}
