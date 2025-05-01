package container

import (
	"errors"
	"sync"

	"github.com/26huitailang/yogo/framework"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider framework.ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// YogoContainer 是服务容器的具体实现
type YogoContainer struct {
	Container
	providers map[string]framework.ServiceProvider
	instances map[string]interface{}
	lock      sync.RWMutex
}

// NewContainer 创建一个新的容器实例
func NewContainer() *YogoContainer {
	return &YogoContainer{
		providers: map[string]framework.ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (c *YogoContainer) Bind(provider framework.ServiceProvider) error {
	c.lock.Lock()
	key := provider.Name()
	c.providers[key] = provider
	c.lock.Unlock()

	if !provider.IsDefer() {
		if err := provider.Boot(c); err != nil {
			return err
		}
		params := provider.Params(c)
		method := provider.Register(c)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		c.instances[key] = instance
	}
	return nil
}

func (c *YogoContainer) IsBind(key string) bool {
	return c.findServiceProvider(key) != nil
}

func (c *YogoContainer) Make(key string) (interface{}, error) {
	return c.make(key, nil, false)
}

func (c *YogoContainer) MustMake(key string) interface{} {
	serv, err := c.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (c *YogoContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.make(key, params, true)
}

func (c *YogoContainer) findServiceProvider(key string) framework.ServiceProvider {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if sp, ok := c.providers[key]; ok {
		return sp
	}
	return nil
}

func (c *YogoContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	// 查找服务提供者
	sp := c.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return c.newInstance(sp, params)
	}

	// 如果容器中已经实例化了，直接返回
	if ins, ok := c.instances[key]; ok {
		return ins, nil
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := c.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	c.instances[key] = inst
	return inst, nil
}

func (c *YogoContainer) newInstance(sp framework.ServiceProvider, params []interface{}) (interface{}, error) {
	// 启动
	if err := sp.Boot(c); err != nil {
		return nil, err
	}

	if params == nil {
		params = sp.Params(c)
	}

	method := sp.Register(c)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, nil
}
