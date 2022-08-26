package http

import (
	"github.com/26huitailang/yogo/app/http/module/demo"
	"github.com/26huitailang/yogo/framework/gin"
	"github.com/26huitailang/yogo/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	demo.Register(r)
}
