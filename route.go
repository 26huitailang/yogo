package main

import (
	"github.com/26huitailang/yogo/framework/gin"
	"github.com/26huitailang/yogo/framework/middleware"
)

func registerRouter(core *gin.Engine) {
	core.Use(
		gin.Recovery(),
		middleware.Cost(),
	)
	core.GET("/user/login", UserLoginController)

	subjectApi := core.Group("/subject")
	subjectApi.Use(middleware.Test2())
	{
		subjectApi.GET("/list/all", SubjectListController)
		subjectApi.GET("/:id", SubjectGetController)
		subjectApi.DELETE("/:id", SubjectDeleteController)
		subjectApi.PUT("/:id", SubjectUpdateController)
	}
}
