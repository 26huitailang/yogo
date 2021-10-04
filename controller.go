package main

import (
	"fmt"
	"github.com/26huitailang/yogo/framework/gin"
	"time"
)

func UserLoginController(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.ISetOkStatus().IJson("ok, UserLoginController")
}

func SubjectListController(c *gin.Context) {
	type subject struct {
		Name string
		Id   int
	}
	c.ISetOkStatus().IJson([]*subject{{"hello", 1}, {"world", 2}})
}

func SubjectGetController(c *gin.Context) {
	c.ISetOkStatus().IJson(fmt.Sprintf("ok, SubjectGetController: %s", c.YogoParam("id")))
}

func SubjectUpdateController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectGetController")
}

func SubjectDeleteController(c *gin.Context) {
	c.ISetOkStatus().IJson("ok, SubjectGetController")
}
