package http

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetContainer(container)
	r.Use(gin.Recovery())
	Routes(r)
	return r, nil
}
