package main

import (
	"github.com/26huitailang/yogo/app/console"
	"github.com/26huitailang/yogo/app/http"
	"github.com/26huitailang/yogo/app/provider/demo"
	"github.com/26huitailang/yogo/framework/container"
	"github.com/26huitailang/yogo/framework/provider/app"
	"github.com/26huitailang/yogo/framework/provider/cache"
	"github.com/26huitailang/yogo/framework/provider/config"
	"github.com/26huitailang/yogo/framework/provider/distributed"
	"github.com/26huitailang/yogo/framework/provider/env"
	"github.com/26huitailang/yogo/framework/provider/kernel"
	"github.com/26huitailang/yogo/framework/provider/log"
	"github.com/26huitailang/yogo/framework/provider/orm"
	"github.com/26huitailang/yogo/framework/provider/redis"
	"github.com/26huitailang/yogo/framework/provider/ssh"
)

func main() {
	c := container.NewContainer()
	c.Bind(&app.YogoAppProvider{})
	// 其他服务提供者绑定
	c.Bind(&distributed.LocalDistributedProvider{})
	c.Bind(&env.YogoEnvProvider{})
	c.Bind(&config.YogoConfigProvider{})
	c.Bind(&log.YogoLogServiceProvider{})
	c.Bind(&orm.GormProvider{})
	c.Bind(&redis.RedisProvider{})
	c.Bind(&cache.YogoCacheProvider{})
	c.Bind(&ssh.SSHProvider{})
	c.Bind(&demo.DemoProvider{})

	if engine, err := http.NewHttpEngine(c); err == nil {
		c.Bind(&kernel.YogoKernelProvider{HttpEngine: engine})
	}
	console.RunCommand(c)
}
