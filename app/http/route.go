package http

import (
	"github.com/26huitailang/yogo/app/http/module/demo"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
	ginSwagger "github.com/26huitailang/yogo/framework/middleware/gin-swagger"
	"github.com/26huitailang/yogo/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/26huitailang/yogo/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	demo.Register(r)
}
