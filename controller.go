package main

import (
	"fmt"
	"github.com/26huitailang/yogo/framework"
)

func UserLoginController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, UserLoginController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	type subject struct {
		Name string
		Id   int
	}
	c.SetOkStatus().Json([]*subject{{"hello", 1}, {"world", 2}})
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.SetOkStatus().Json(fmt.Sprintf("ok, SubjectGetController: %s", c.Param("id")))
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectGetController")
	return nil
}

func SubjectDeleteController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectGetController")
	return nil
}
