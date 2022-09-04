package contract

import (
	"fmt"

	"github.com/26huitailang/yogo/framework"
	"github.com/go-redis/redis/v8"
)

const RedisKey = "yogo:redis"

type RedisOption func(container framework.Container, config *RedisConfig) error
type RedisService interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}
type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
