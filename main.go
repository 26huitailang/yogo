package main

import (
	"net/http"

	"github.com/26huitailang/yogo/framework"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
