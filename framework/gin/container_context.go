package gin

import (
	"github.com/26huitailang/yogo/framework/container"
	"github.com/gin-gonic/gin"
)

// ContainerContext 是 gin.Context 的包装器
type ContainerContext struct {
	*gin.Context
	container *container.YogoContainer
}

// NewContainerContext 创建一个新的 ContainerContext
func NewContainerContext(c *gin.Context, container *container.YogoContainer) *ContainerContext {
	return &ContainerContext{
		Context:   c,
		container: container,
	}
}

// Make 从容器中获取服务
func (c *ContainerContext) Make(key string) (interface{}, error) {
	return c.container.Make(key)
}

// MustMake 从容器中获取服务，如果服务不存在则 panic
func (c *ContainerContext) MustMake(key string) interface{} {
	return c.container.MustMake(key)
}

// MakeNew 创建服务实例，不使用单例模式
func (c *ContainerContext) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.container.MakeNew(key, params)
}

// IsBind 判断服务是否已绑定
func (c *ContainerContext) IsBind(key string) bool {
	return c.container.IsBind(key)
}

// GetContainer 获取容器实例
func (c *ContainerContext) GetContainer() *container.YogoContainer {
	return c.container
}
