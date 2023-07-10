package http

import (
	"errors"
	"time"

	"github.com/26huitailang/yogo/app/http/module/demo"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/gin"
	ginSwagger "github.com/26huitailang/yogo/framework/middleware/gin-swagger"
	limiter "github.com/26huitailang/yogo/framework/middleware/gin-limiter"
	"github.com/26huitailang/yogo/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/26huitailang/yogo/framework/middleware/static"
)

func Routes(r *gin.Engine) {
	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	lm := limiter.NewRateLimiter(time.Second, 5, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-USER-TOKEN")
		if key != "" {
			return key, nil
		}
		return "", errors.New("User is not authorized")
	})

	r.Use(lm.Middleware())
	demo.Register(r)
}
