package main

import (
	"net/http"

	"github.com/26huitailang/yogo/framework"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
