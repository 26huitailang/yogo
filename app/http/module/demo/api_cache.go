package demo

import (
	"time"

	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
)

// DemoCache cache的简单例子
func (api *DemoApi) DemoCache(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)
	// 初始化cache服务
	cacheService := c.MustMake(contract.CacheKey).(contract.CacheService)
	// 设置key为foo
	err := cacheService.Set(c, "foo", "bar", 1*time.Hour)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// 获取key为foo
	val, err := cacheService.Get(c, "foo")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "cache get", map[string]interface{}{
		"val": val,
	})
	// 删除key为foo
	if err := cacheService.Del(c, "foo"); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, "ok")
}
