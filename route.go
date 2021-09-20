package main

import "github.com/26huitailang/yogo/framework"

func registerRouter(core *framework.Core) {
	core.Get("/user/login", UserLoginController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.Get("/list/all", SubjectListController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Delete("/:id", SubjectDeleteController)
		subjectApi.Put("/:id", SubjectUpdateController)
	}
}
