package http

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.ContainerEngine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.NewContainerEngine()
	r.SetContainer(container)
	r.Use(gin.Recovery())
	Routes(r)
	return r, nil
}
