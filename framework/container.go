package framework

import (
	"errors"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。  // 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// YogoContainer 是服务容器的具体实现
type YogoContainer struct {
	Container
	providers map[string]ServiceProvider
	instances map[string]interface{}
	lock      sync.RWMutex
}

func NewYogoContainer() *YogoContainer {
	return &YogoContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (c *YogoContainer) Bind(provider ServiceProvider) error {
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

// Make 方式调用内部的 make实现
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

func (c *YogoContainer) NameList() []string {
	ret := []string{}
	for _, provider := range c.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}

func (c *YogoContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	sp := c.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return c.newInstance(sp, params)
	}

	if ins, ok := c.instances[key]; ok {
		return ins, nil
	}

	ins, err := c.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	c.instances[key] = ins
	return ins, nil
}

func (c *YogoContainer) findServiceProvider(key string) ServiceProvider {
	sp, ok := c.providers[key]
	if !ok {
		return nil
	}
	return sp
}

func (c *YogoContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(c); err != nil {
		return nil, err
	}

	if params == nil {
		params = provider.Params(c)
	}
	method := provider.Register(c)
	instance, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return instance, nil
}
