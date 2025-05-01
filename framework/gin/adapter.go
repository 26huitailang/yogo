package gin

import (
	"github.com/gin-gonic/gin"
)

// WrapHandler 将 ContainerContext 处理函数转换为 gin.HandlerFunc
func WrapHandler(handler func(*ContainerContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 gin.Context 中获取 ContainerContext
		if ctx, exists := c.Keys["container"]; exists {
			if containerCtx, ok := ctx.(*ContainerContext); ok {
				handler(containerCtx)
				return
			}
		}
		// 如果没有找到 ContainerContext，创建一个新的
		if engine, exists := c.Keys["engine"]; exists {
			if containerEngine, ok := engine.(*ContainerEngine); ok {
				containerCtx := NewContainerContext(c, containerEngine.container)
				handler(containerCtx)
				return
			}
		}
		// 如果都没有找到，使用原始的 gin.Context
		c.Next()
	}
}

// WrapHandlers 将多个 ContainerContext 处理函数转换为 gin.HandlersChain
func WrapHandlers(handlers ...func(*ContainerContext)) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = WrapHandler(h)
	}
	return funcs
}
