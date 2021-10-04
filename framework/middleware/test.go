package middleware

import (
	"fmt"

	"github.com/26huitailang/yogo/framework/gin"
)

func Test1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test1")
		c.Next()
		fmt.Println("middleware post test1")
	}
}

func Test2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware post test2")
	}
}
