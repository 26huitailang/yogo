package main

import (
	"net/http"

	"github.com/26huitailang/yogo/framework"
)

func FooControllerHandler(ctx *framework.Context) error {
	return ctx.Json(http.StatusOK, map[string]interface{}{
		"code": 0,
	})
}
