package kernel

import (
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
	"net/http"
)

type YogoKernelService struct {
	contract.Kernel
	engine *gin.Engine
}

func NewYogoKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &YogoKernelService{engine: httpEngine}, nil
}

func (s *YogoKernelService) HttpEngine() http.Handler {
	return s.engine
}
