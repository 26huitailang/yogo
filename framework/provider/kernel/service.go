package kernel

import (
	"net/http"

	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
)

type YogoKernelService struct {
	contract.Kernel
	engine *gin.ContainerEngine
}

func NewYogoKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.ContainerEngine)
	return &YogoKernelService{engine: httpEngine}, nil
}

func (s *YogoKernelService) HttpEngine() http.Handler {
	return s.engine
}
