package main

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/middleware"
)

func registerRouter(core *framework.Core) {
	core.Use(
		middleware.Recovery(),
		middleware.Cost(),
	)
	core.Get("/user/login", UserLoginController)

	subjectApi := core.Group("/subject")
	subjectApi.Use(middleware.Test2())
	{
		subjectApi.Get("/list/all", SubjectListController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Delete("/:id", SubjectDeleteController)
		subjectApi.Put("/:id", SubjectUpdateController)
	}
}
