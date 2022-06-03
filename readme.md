
# 服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

> 服务的使用基本按照以下顺序

1. 启用服务, 在 zapp.NewApp 提供选项, 一般为 service.WithService()
2. 注入服务, 这一步根据不同的服务有不同的方法, 一般为 service.RegistryXXX(...)
3. 启动app

# 以api服务示例

```go
app := zapp.NewApp("test",
    api.WithService(), // 启用api服务
)

// 注册路由
api.RegistryRouter(func(c core.IComponent, router api.Party) {
    router.Get("/", api.Wrap(func(ctx *api.Context) interface{} {
        return "hello"
    }))
})

// 运行
app.Run()
```
