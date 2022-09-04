package main

import (
	"github.com/26huitailang/yogo/app/console"
	"github.com/26huitailang/yogo/app/http"
	"github.com/26huitailang/yogo/app/provider/demo"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/provider/app"
	"github.com/26huitailang/yogo/framework/provider/config"
	"github.com/26huitailang/yogo/framework/provider/distributed"
	"github.com/26huitailang/yogo/framework/provider/env"
	"github.com/26huitailang/yogo/framework/provider/kernel"
	"github.com/26huitailang/yogo/framework/provider/log"
	"github.com/26huitailang/yogo/framework/provider/orm"
)

func main() {
	container := framework.NewYogoContainer()
	container.Bind(&app.YogoAppProvider{})
	// 其他服务提供者绑定
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&env.YogoEnvProvider{})
	container.Bind(&config.YogoConfigProvider{})
	container.Bind(&log.YogoLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&demo.DemoProvider{})

	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.YogoKernelProvider{HttpEngine: engine})
	}
	console.RunCommand(container)
}
