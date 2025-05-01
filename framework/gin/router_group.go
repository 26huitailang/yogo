package gin

import (
	"github.com/gin-gonic/gin"
)

// ContainerRouterGroup 是 gin.RouterGroup 的包装器
type ContainerRouterGroup struct {
	*gin.RouterGroup
}

// GET 包装 gin.RouterGroup 的 GET 方法
func (group *ContainerRouterGroup) GET(path string, handlers ...func(*ContainerContext)) *ContainerRouterGroup {
	group.RouterGroup.GET(path, WrapHandlers(handlers...)...)
	return group
}

// POST 包装 gin.RouterGroup 的 POST 方法
func (group *ContainerRouterGroup) POST(path string, handlers ...func(*ContainerContext)) *ContainerRouterGroup {
	group.RouterGroup.POST(path, WrapHandlers(handlers...)...)
	return group
}

// Group 包装 gin.RouterGroup 的 Group 方法
func (group *ContainerRouterGroup) Group(path string, handlers ...func(*ContainerContext)) *ContainerRouterGroup {
	return &ContainerRouterGroup{
		RouterGroup: group.RouterGroup.Group(path, WrapHandlers(handlers...)...),
	}
}

// Use 包装 gin.RouterGroup 的 Use 方法
func (group *ContainerRouterGroup) Use(handlers ...func(*ContainerContext)) *ContainerRouterGroup {
	group.RouterGroup.Use(WrapHandlers(handlers...)...)
	return group
}
