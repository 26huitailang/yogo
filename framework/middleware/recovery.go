package middleware

import (
	"github.com/26huitailang/yogo/framework"
	"net/http"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(http.StatusInternalServerError)
				c.Json(err)
			}
		}()

		c.Next()

		return nil
	}
}
