目前该项目代码侵入修改了gin的代码，我想保持现有service的思路，但是不想侵入源码修改

## 重构方案

为了避免修改gin源码，同时保持现有的service功能，我们可以采用包装器模式（Wrapper Pattern）来重构代码。主要步骤如下：

### 1. 创建容器包

首先创建独立的容器包，将容器相关的代码从gin中分离出来：

```go
// framework/container/container.go
package container

type Container interface {
    Bind(provider framework.ServiceProvider) error
    IsBind(key string) bool
    Make(key string) (interface{}, error)
    MustMake(key string) interface{}
    MakeNew(key string, params []interface{}) (interface{}, error)
}

type YogoContainer struct {
    Container
    providers map[string]framework.ServiceProvider
    instances map[string]interface{}
    lock      sync.RWMutex
}
```

### 2. 创建gin的包装器

创建一个包装gin.Engine的ContainerEngine，以及包装gin.Context的ContainerContext：

```go
// framework/gin/wrapper.go
package gin

type ContainerEngine struct {
    *origingin.Engine
    container framework.Container
}

type ContainerContext struct {
    *origingin.Context
    container framework.Container
}

type ContainerHandlerFunc func(*ContainerContext)
```

主要功能：
- ContainerEngine 包装了gin.Engine，添加了容器支持
- ContainerContext 包装了gin.Context，提供了容器相关的方法
- 所有的HTTP方法(GET/POST等)都被重写为使用ContainerContext

### 3. 修改现有代码

需要修改的地方：

1. HTTP引擎初始化：
```go
// app/http/kernel.go
func NewHttpEngine(container framework.Container) (*gin.ContainerEngine, error) {
    r := gin.NewContainerEngine()
    r.SetContainer(container)
    // ...
    return r, nil
}
```

2. 路由定义：
```go
// app/http/route.go
func Routes(r *gin.ContainerEngine) {
    container := r.GetContainer()
    // ...
}
```

3. 控制器方法：
```go
// app/http/module/demo/api.go
func (api *DemoApi) Demo(c *gin.ContainerContext) {
    c.JSON(200, "this is demo")
}
```

## 优点

1. 不再需要修改gin的源码，更容易升级gin版本
2. 保持了原有的服务容器功能
3. 代码结构更清晰，更容易维护
4. 通过组合而不是继承来扩展功能，符合Go的设计理念

## 注意事项

1. 需要修改所有使用原来gin.Engine和gin.Context的代码
2. 中间件可能需要适配新的Context类型
3. 第三方包如果依赖gin.Context，需要在使用时进行类型转换
4. 确保所有的路由处理函数都使用新的ContainerContext

## 迁移步骤

1. 添加gin原版依赖：
```bash
go get -u github.com/gin-gonic/gin
```

2. 创建新的包装器文件和容器包
3. 修改现有的控制器代码，将Context参数改为ContainerContext
4. 修改路由注册代码，使用新的ContainerEngine
5. 修改中间件，确保它们能够正确处理ContainerContext
6. 进行完整的测试，确保功能正常