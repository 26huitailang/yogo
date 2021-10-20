package main

import (
	"github.com/26huitailang/yogo/framework/gin"
)

func registerRouter(core *gin.Engine) {
	core.GET("/user/login", UserLoginController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.GET("/list/all", SubjectListController)
		subjectApi.GET("/:id", SubjectGetController)
		subjectApi.DELETE("/:id", SubjectDeleteController)
		subjectApi.PUT("/:id", SubjectUpdateController)
	}
}
