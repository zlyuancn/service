
# nsq消费服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
nsq_consume.RegistryService()               # 注册服务
nsq_consume.WithNsqConsumeService()         # 启用服务
nsq_consume.RegistryNsqConsumeHandler(...)      # 服务注入(注册消费处理handler)
```

# 示例

```go
package main

import (
	nsq_consume "github.com/zly-app/service/nsq-consume"
	"github.com/zly-app/zapp"
)

func main() {
	// 注册nsq消费服务
	nsq_consume.RegistryService()
	// 启用nsq消费服务
	app := zapp.NewApp("test", nsq_consume.WithNsqConsumeService())
	// 注册handler
	nsq_consume.RegistryNsqConsumeHandler(app, "test", "c2", func(ctx *nsq_consume.Context) error {
		ctx.Info("数据", string(ctx.Body))
		return nil
	})
	// 运行
	app.Run()
}
```

# 配置

> 默认服务类型为 `nsq-consume`

```toml
[services.nsq-consume]
# nsqd地址, localhost1:4150,localhost2:4150
NsqdAddress="localhost:4150"
# nsq发现服务地址, 优先级高于NsqdAddress, localhost1:4161,localhost2:4161
NsqLookupdAddress=""
# 验证秘钥
AuthSecret=""
# 心跳间隔(毫秒), 不能超过ReadTimeout
HeartbeatInterval=30000
# 超时(毫秒)
ReadTimeout=30000
# 超时(毫秒)
WriteTimeout=5000
# 超时(毫秒)
DialTimeout=2000
# 线程数, 默认为0表示使用逻辑cpu数量
ThreadCount=0
```
