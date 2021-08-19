# kafka消费服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
kafka_consume.WithService()             # 启用服务
kafka_consume.RegistryHandler(...)      # 服务注入(注册消费处理handler)
```

# 示例

```go
package main

import (
	kafka_consume "github.com/zly-app/service/kafka-consume"
	"github.com/zly-app/zapp"
)

func main() {
	// 启用kafka消费服务
	app := zapp.NewApp("test", kafka_consume.WithService())
	// 注册handler
	kafka_consume.RegistryHandler(app, []string{"test"}, "g3", func(ctx *kafka_consume.Context) error {
		ctx.Info("数据", string(ctx.Value))
		return nil
	})
	// 运行
	app.Run()
}
```

# 配置

> 默认服务类型为 `kafka-consume`, 完整配置说明参考 [Config](./config.go)

```toml
[services.kafka-consume]
Address = "localhost:9092"                # 地址, 多个地址用半角逗号连接
```
