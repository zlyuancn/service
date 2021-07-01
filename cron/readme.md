
# cron服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
cron.WithService()              # 启用服务
cron.RegistryHandler(...)       # 注册handler
cron.RegistryOnceHandler(...)   # 注册一次性handler
cron.RegistryTask(...)          # 注册自定义task
```

# 示例
```go
package main

import (
	"github.com/zly-app/service/cron"
	"github.com/zly-app/zapp"
)

func main() {
	// 启用cron服务
	app := zapp.NewApp("test", cron.WithService())
	// 注册handler
	cron.RegistryHandler(app, "c1", "@every 1s", true, func(ctx cron.IContext) error {
		ctx.Info("触发")
		return nil
	})
	// 运行
	app.Run()
}
```

# 配置

> 这个服务可以不需要配置, 默认服务类型为 `cron`

```toml
[services.cron]
# 线程数, 默认为-1
ThreadCount = -1
# 最大任务队列大小, 默认为10000
MaxTaskQueueSize = 0
```
