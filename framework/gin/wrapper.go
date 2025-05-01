package gin

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/container"
	origingin "github.com/gin-gonic/gin"
)

// ContainerEngine 包装gin.Engine，添加容器支持
type ContainerEngine struct {
	*origingin.Engine
	container framework.Container
}

// NewContainerEngine 创建一个新的Engine实例
func NewContainerEngine() *ContainerEngine {
	engine := &ContainerEngine{
		Engine:    origingin.New(),
		container: container.NewContainer(),
	}
	return engine
}

// DefaultContainerEngine 创建一个带有默认中间件的Engine实例
func DefaultContainerEngine() *ContainerEngine {
	engine := &ContainerEngine{
		Engine:    origingin.Default(),
		container: container.NewContainer(),
	}
	return engine
}

// SetContainer 设置容器
func (engine *ContainerEngine) SetContainer(container framework.Container) {
	engine.container = container
}

// GetContainer 获取容器
func (engine *ContainerEngine) GetContainer() framework.Container {
	return engine.container
}

// Bind 绑定服务提供者
func (engine *ContainerEngine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 检查服务是否已绑定
func (engine *ContainerEngine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

// ContainerContext 包装gin.Context，添加容器支持
type ContainerContext struct {
	*origingin.Context
	container framework.Container
}

// ContainerHandlerFunc 定义处理函数类型
type ContainerHandlerFunc func(*ContainerContext)

// Handle 包装gin的Handle方法，支持容器注入
func (engine *ContainerEngine) Handle(httpMethod, relativePath string, handlers ...ContainerHandlerFunc) {
	ginHandlers := make([]origingin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = func(c *origingin.Context) {
			ctx := &ContainerContext{
				Context:   c,
				container: engine.container,
			}
			handler(ctx)
		}
	}
	engine.Engine.Handle(httpMethod, relativePath, ginHandlers...)
}

// GET 是Handle("GET", path, handle)的快捷方式
func (engine *ContainerEngine) GET(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("GET", relativePath, handlers...)
}

// POST 是Handle("POST", path, handle)的快捷方式
func (engine *ContainerEngine) POST(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("POST", relativePath, handlers...)
}

// PUT 是Handle("PUT", path, handle)的快捷方式
func (engine *ContainerEngine) PUT(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("PUT", relativePath, handlers...)
}

// DELETE 是Handle("DELETE", path, handle)的快捷方式
func (engine *ContainerEngine) DELETE(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("DELETE", relativePath, handlers...)
}

// PATCH 是Handle("PATCH", path, handle)的快捷方式
func (engine *ContainerEngine) PATCH(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("PATCH", relativePath, handlers...)
}

// HEAD 是Handle("HEAD", path, handle)的快捷方式
func (engine *ContainerEngine) HEAD(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("HEAD", relativePath, handlers...)
}

// OPTIONS 是Handle("OPTIONS", path, handle)的快捷方式
func (engine *ContainerEngine) OPTIONS(relativePath string, handlers ...ContainerHandlerFunc) {
	engine.Handle("OPTIONS", relativePath, handlers...)
}

// Make 从容器中获取服务
func (c *ContainerContext) Make(key string) (interface{}, error) {
	return c.container.Make(key)
}

// MustMake 从容器中获取服务，如果不存在则panic
func (c *ContainerContext) MustMake(key string) interface{} {
	return c.container.MustMake(key)
}

// MakeNew 从容器中获取服务的新实例
func (c *ContainerContext) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.container.MakeNew(key, params)
}
