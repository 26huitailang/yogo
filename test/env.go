package tests

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/provider/app"
	"github.com/26huitailang/yogo/framework/provider/env"
)

const (
	BasePath = "/Users/26huitailang/GolangProjects/yogo-demo/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewYogoContainer()
	// 绑定App服务提供者
	container.Bind(&app.YogoAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.YogoTestingEnvProvider{})
	return container
}
