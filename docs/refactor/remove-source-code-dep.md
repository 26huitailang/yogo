目前该项目代码侵入修改了gin的代码，我想保持现有service的思路，但是不想侵入源码修改

## 重构方案

为了避免修改gin源码，同时保持现有的service功能，我们采用包装器模式（Wrapper Pattern）来重构代码。主要步骤如下：

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
// framework/gin/container_engine.go
type ContainerEngine struct {
    *gin.Engine
    container framework.Container // 使用接口而不是具体类型
}

// framework/gin/container_context.go
type ContainerContext struct {
    *gin.Context
    container framework.Container // 使用接口而不是具体类型
}
```

### 3. 实现的主要功能

1. ContainerEngine 包装了 gin.Engine：
   - 添加了容器支持
   - 实现了 GetContainer() 和 SetContainer() 方法
   - 包装了所有主要的 HTTP 方法（GET、POST 等）
   - 实现了路由组的包装

2. ContainerContext 包装了 gin.Context：
   - 提供了容器相关的方法（Make、MustMake 等）
   - 实现了服务获取的便捷方法
   - 保持了与原有 gin.Context 的兼容性

3. 适配器层：
   - 实现了 WrapHandler 和 WrapHandlers 函数
   - 处理了中间件的适配
   - 保证了与现有代码的兼容性

### 4. 重构进展

已完成的工作：
1. 移除了对 gin 源码的直接修改
2. 实现了基于接口的依赖注入
3. 完成了主要包装器的实现
4. 修复了编译错误
5. 实现了更好的解耦

主要改进：
1. 使用 `framework.Container` 接口替代具体的 `*container.YogoContainer` 类型
2. 重构了容器的创建方式，使用 `container.NewContainer()`
3. 更新了所有相关方法签名，使用接口而不是具体类型
4. 实现了更好的依赖倒置原则

### 5. 优点

1. 不再需要修改gin的源码，更容易升级gin版本
2. 保持了原有的服务容器功能
3. 代码结构更清晰，更容易维护
4. 通过组合而不是继承来扩展功能，符合Go的设计理念
5. 更容易进行单元测试
6. 更容易扩展和替换实现
7. 遵循了依赖倒置原则

### 6. 注意事项

1. 需要修改所有使用原来gin.Engine和gin.Context的代码
2. 中间件可能需要适配新的Context类型
3. 第三方包如果依赖gin.Context，需要在使用时进行类型转换
4. 确保所有的路由处理函数都使用新的ContainerContext

### 7. 后续工作

1. 完善错误处理：
   - 添加更详细的错误信息
   - 实现统一的错误处理机制

2. 优化中间件：
   - 创建专门的中间件适配层
   - 实现常用中间件的转换函数

3. 完善文档：
   - 为核心接口添加详细文档
   - 编写使用示例和最佳实践

4. 添加测试：
   - 编写单元测试
   - 添加集成测试

5. 性能优化：
   - 优化容器的实例管理
   - 添加缓存机制

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
