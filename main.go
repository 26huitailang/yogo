package main

import (
	"github.com/26huitailang/yogo/app/console"
	"github.com/26huitailang/yogo/app/http"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/provider/app"
	"github.com/26huitailang/yogo/framework/provider/distributed"
	"github.com/26huitailang/yogo/framework/provider/env"
	"github.com/26huitailang/yogo/framework/provider/kernel"
)

func main() {
	container := framework.NewYogoContainer()
	container.Bind(&app.YogoAppProvider{})
	// 其他服务提供者绑定
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&env.YogoEnvProvider{})

	if engine, err := http.NewHttpEngine(); err == nil {
		container.Bind(&kernel.YogoKernelProvider{HttpEngine: engine})
	}
	console.RunCommand(container)
}
