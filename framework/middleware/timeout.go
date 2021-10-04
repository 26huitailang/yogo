package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/26huitailang/yogo/framework/gin"
)

func TimeoutHandler(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前预操作：初始化超时 context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.ISetStatus(http.StatusInternalServerError).IJson("time out")
			log.Println(p)
		case <-finish:
			log.Println("finish")
		case <-durationCtx.Done():
			c.ISetStatus(http.StatusInternalServerError).IJson("time out")
		}
	}
}
