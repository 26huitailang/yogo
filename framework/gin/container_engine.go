package gin

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/container"
	"github.com/gin-gonic/gin"
)

// ContainerEngine 是一个继承 gin.Engine 的结构体
type ContainerEngine struct {
	*gin.Engine
	container *container.YogoContainer // 服务容器
}

// New 创建一个新的 ContainerEngine 实例
func New() *ContainerEngine {
	engine := &ContainerEngine{
		Engine:    gin.New(),
		container: container.NewContainer(),
	}
	// 设置 engine 的 HandleContext
	engine.Engine.HandleContext = engine.HandleContext
	return engine
}

// HandleContext 包装 gin 的上下文处理
func (engine *ContainerEngine) HandleContext(c *gin.Context) {
	ctx := NewContainerContext(c, engine.container)
	c.Keys = make(map[string]interface{})
	c.Keys["container"] = ctx
	engine.Engine.HandleContext(c)
}

// Bind 绑定服务提供者
func (engine *ContainerEngine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 判断服务是否已绑定
func (engine *ContainerEngine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

// Make 获取服务实例
func (engine *ContainerEngine) Make(key string) (interface{}, error) {
	return engine.container.Make(key)
}

// MustMake 获取服务实例，如果服务不存在则 panic
func (engine *ContainerEngine) MustMake(key string) interface{} {
	return engine.container.MustMake(key)
}

// MakeNew 创建服务实例，不使用单例模式
func (engine *ContainerEngine) MakeNew(key string, params []interface{}) (interface{}, error) {
	return engine.container.MakeNew(key, params)
}
