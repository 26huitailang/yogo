package main

import (
	"github.com/26huitailang/yogo/framework"
)

func UserLoginController(c *framework.Context) error {
	c.Json(200, "ok, UserLoginController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	type subject struct {
		Name string
		Id   int
	}
	c.Json(200, []*subject{{"hello", 1}, {"world", 2}})
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.Json(200, "ok, SubjectGetController")
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.Json(200, "ok, SubjectGetController")
	return nil
}

func SubjectDeleteController(c *framework.Context) error {
	c.Json(200, "ok, SubjectGetController")
	return nil
}
