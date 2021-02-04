
# 服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

> 服务的使用基本按照以下顺序

1. 注册服务, 一般为 serviceType.RegistryService()
2. 启用服务, 在 zapp.NewApp 提供选项, 一般为 serviceType.WithXXXService()
3. 注入服务, 这一步根据不同的服务有不同的方法, 一般为 serviceType.RegistryXXX(...)
4. 启动app

# 以api服务示例

```go
// 注册api服务
api.RegistryService()
// ... 注册其他服务, 一般为 serviceType.RegistryService()

app := zapp.NewApp("test",
    api.WithApiService(), // 启用api服务
    // ... 启用其它服务, 一般为 serviceType.WithXXXService()
)

// 注册路由
api.RegistryApiRouter(app, func(c core.IComponent, router api.Party) {
    router.Get("/", api.Wrap(func(ctx *api.Context) interface{} {
        return "hello"
    }))
})

// ... 其它服务的注入, 一般为 serviceType.RegistryXXX(...)

// 运行
app.Run()
```
