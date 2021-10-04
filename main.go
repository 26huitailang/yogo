package main

import (
	"context"
	"github.com/26huitailang/yogo/framework/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := gin.New()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	server.RegisterOnShutdown(func() {
		log.Println("register shutdown")
	})
	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
