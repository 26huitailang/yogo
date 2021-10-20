package main

import (
	"fmt"
	"github.com/26huitailang/yogo/framework/gin"
	"github.com/26huitailang/yogo/provider/demo"
	"time"
)

func UserLoginController(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.ISetOkStatus().IJson("ok, UserLoginController")
}

func SubjectListController(c *gin.Context) {
	demoService := c.MustMake(demo.Key).(demo.Service)
	foo := demoService.GetFoo()
	c.ISetOkStatus().IJson(foo)
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
