package http

import "github.com/26huitailang/yogo/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	Routes(r)
	return r, nil
}
