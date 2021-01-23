
# 服务插件

> 提供用于 https://github.com/zly-app/zapp 的服务插件

# 以api服务插件示例

```go
// 注册api服务
api.RegistryService()
// ... 注册其他服务, 一般为 type.RegistryService()

app := zapp.NewApp("test",
    api.WithApiService(), // 启用api服务
    // ... 启用其它服务, 一般为 type.WithXXXService()
)

// 注册路由
api.RegistryApiRouter(app, func(c core.IComponent, router api.Party) {
    router.Get("/", api.Wrap(func(ctx *api.Context) interface{} {
        return "hello"
    }))
})

// ... 其它服务的注入, 一般为 type.RegistryXXX(...)

// 运行
app.Run()
```
