package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/26huitailang/yogo/framework"
)

func FooControllerHandler(ctx *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		time.Sleep(6 * time.Second)
		ctx.Json(200, "OK")
		finish <- struct{}{}
	}()
	select {
	case p := <-panicChan:
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		log.Println(p)
		ctx.Json(http.StatusInternalServerError, "panic")
	case <-finish:
		log.Println("finish")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.Json(http.StatusInternalServerError, "time out")
		ctx.SetHasTimeout()
	}
	return nil
}

func UserLoginController(c *framework.Context) error {
	c.Json(200, "OK")
	return nil
}

func SubjectListController(c *framework.Context) error {
	type subject struct {
		Name string
		Id   int
	}
	c.Json(200, []*subject{{"hello", 1}, {"world", 2}})
	return nil
}
