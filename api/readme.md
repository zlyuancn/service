
# api服务插件

> 提供用于 https://github.com/zly-app/zapp 的服务插件

# 说明

```text
api.RegistryService()           # 注册服务
api.WithApiService()            # 启用服务
api.RegistryApiRouter(...)      # 服务注入(注册路由)
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
