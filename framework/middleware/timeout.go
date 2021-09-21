package middleware

import (
	"context"
	"log"
	"time"

	"github.com/26huitailang/yogo/framework"
)

func TimeoutHandler(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
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
			log.Println(p)
			c.Json(500, "time out")
		case <-finish:
			log.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.Json(500, "time out")
		}
		return nil
	}
}
