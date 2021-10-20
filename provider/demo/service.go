package demo

import (
	"fmt"
	"github.com/26huitailang/yogo/framework"
)

type DemoService struct {
	Service
	c framework.Container
}

// 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	fmt.Println("new demo service")
	// 返回实例
	return &DemoService{c: c}, nil
}

func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "I'm foo",
	}
}
