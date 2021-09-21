package middleware

import (
	"log"
	"time"

	"github.com/26huitailang/yogo/framework"
)

func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		uri := c.GetRequest().RequestURI
		method := c.GetRequest().Method
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("%s %s: %s", method, uri, cost)
		return nil
	}
}
