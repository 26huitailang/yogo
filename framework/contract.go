package framework

// NewInstance 定义了创建新实例的函数类型
type NewInstance func(...interface{}) (interface{}, error)

// Container 定义服务容器的接口
type Container interface {
	// Bind 绑定一个服务提供者
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool
	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// ServiceProvider 定义了服务提供者的接口
type ServiceProvider interface {
	// Name 返回服务提供者的凭证
	Name() string
	// Register 注册一个服务实例
	Register(Container) NewInstance
	// Boot 启动的时候判断是否要实例化
	Boot(Container) error
	// IsDefer 是否延迟实例化
	IsDefer() bool
	// Params 获取初始化参数
	Params(Container) []interface{}
}
