# grpc服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
grpc.WithService()              # 启用服务
grpc.RegistryServerObject(...)  # 服务注入(注册服务实体)
```

# 示例

+ [grpc服务端](./example/grpc-s)
+ [grpc客户端](https://github.com/zly-app/component/tree/master/grpc-client/example/grpc-c)

# 配置

> 默认服务类型为 `grpc`

```toml
[services.grpc]
# bind地址
Bind = "localhost:3000"
```
