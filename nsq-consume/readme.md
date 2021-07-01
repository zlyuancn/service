
# nsq消费服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
nsq_consume.WithService()           # 启用服务
nsq_consume.RegistryHandler(...)    # 服务注入(注册消费处理handler)
```

# 示例

```go
package main

import (
	nsq_consume "github.com/zly-app/service/nsq-consume"
	"github.com/zly-app/zapp"
)

func main() {
	// 启用nsq消费服务
	app := zapp.NewApp("test", nsq_consume.WithService())
	// 注册handler
	nsq_consume.RegistryHandler(app, "test", "c2", func(ctx *nsq_consume.Context) error {
		ctx.Info("数据", string(ctx.Body))
		return nil
	})
	// 运行
	app.Run()
}
```

# 配置

> 默认服务类型为 `nsq-consume`, 完整配置说明参考 [Config](./config.go)

```toml
[services.nsq-consume]
# nsqd地址, localhost1:4150,localhost2:4150
NsqdAddress="localhost:4150"
# nsq发现服务地址, 优先级高于NsqdAddress, localhost1:4161,localhost2:4161
NsqLookupdAddress="localhost:4161"
```
