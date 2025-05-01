package gin

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/gin-gonic/gin"
)

// ContainerEngine 是一个继承 gin.Engine 的结构体
type ContainerEngine struct {
	*gin.Engine
	container framework.Container // 服务容器
}

// New 创建一个新的 ContainerEngine 实例
func New() *ContainerEngine {
	engine := &ContainerEngine{
		Engine: gin.New(),
	}
	// 使用自定义的中间件来处理上下文
	engine.Engine.Use(func(c *gin.Context) {
		ctx := NewContainerContext(c, engine.container)
		c.Keys = make(map[string]interface{})
		c.Keys["container"] = ctx
		c.Keys["engine"] = engine
		c.Next()
	})
	return engine
}

// GET 包装 gin.Engine 的 GET 方法
func (engine *ContainerEngine) GET(path string, handlers ...func(*ContainerContext)) *ContainerEngine {
	engine.Engine.GET(path, WrapHandlers(handlers...)...)
	return engine
}

// POST 包装 gin.Engine 的 POST 方法
func (engine *ContainerEngine) POST(path string, handlers ...func(*ContainerContext)) *ContainerEngine {
	engine.Engine.POST(path, WrapHandlers(handlers...)...)
	return engine
}

// Group 包装 gin.Engine 的 Group 方法
func (engine *ContainerEngine) Group(path string, handlers ...func(*ContainerContext)) *ContainerRouterGroup {
	return &ContainerRouterGroup{
		RouterGroup: engine.Engine.Group(path, WrapHandlers(handlers...)...),
	}
}

// Use 包装 gin.Engine 的 Use 方法
func (engine *ContainerEngine) Use(handlers ...func(*ContainerContext)) *ContainerEngine {
	engine.Engine.Use(WrapHandlers(handlers...)...)
	return engine
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

// NewContainerEngine 是 New 的别名，用于向后兼容
func NewContainerEngine() *ContainerEngine {
	return New()
}

// Recovery 返回一个 Recovery 中间件
func Recovery() func(*ContainerContext) {
	return func(c *ContainerContext) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// GetContainer 返回服务容器
func (engine *ContainerEngine) GetContainer() framework.Container {
	return engine.container
}

// SetContainer 设置服务容器
func (engine *ContainerEngine) SetContainer(container framework.Container) {
	engine.container = container
}
