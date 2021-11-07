package demo

import (
	"github.com/26huitailang/yogo/framework"
)

type Service struct {
	c framework.Container
}

// 初始化实例的方法
func NewService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)
	return &Service{c: c}, nil
}

func (s *Service) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}
}
