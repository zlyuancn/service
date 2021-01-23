
# api服务插件

> 提供用于 https://github.com/zly-app/zapp 的服务插件

# 示例

```go
// 注册api服务
api.RegistryService()
// 启用api服务
app := zapp.NewApp("test", api.WithApiService())
// 注册路由
api.RegistryApiRouter(app, func(c core.IComponent, router api.Party) {
    router.Get("/", api.Wrap(func(ctx *api.Context) interface{} {
        return "hello"
    }))
})
// 运行
app.Run()
```

# 配置

> 默认服务类型为 `api`

```toml
[services.api]
# bind地址
Bind=":8080"
# 适配nginx的Forwarded获取ip, 优先级高于nginx的Real
IPWithNginxForwarded=true
# 适配nginx的Real获取ip, 优先级高于sock连接的ip
IPWithNginxReal=true
# 在开发环境中显示api结果
ShowApiResultInDevelop=true
# 在生产环境显示详细的错误
ShowDetailedErrorInProduction=false
```
