package middleware

import (
	"log"
	"time"

	"github.com/26huitailang/yogo/framework/gin"
)

func Cost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		method := c.Request.Method
		start := time.Now()

		c.Next()

		end := time.Now()
		cost := end.Sub(start)
		log.Printf("%s %s: %s", method, uri, cost)
	}
}
