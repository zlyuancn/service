
# cron服务插件

> 提供用于 https://github.com/zly-app/zapp 的服务插件

# 说明

> 此组件基于模块 [github.com/zlyuancn/zscheduler](https://github.com/zlyuancn/zscheduler)

```text
cron.RegistryService()             # 注册服务
cron.WithCronService()             # 启用服务
cron.RegistryCronHandler(...)           # 服务注入(注册handler)
cron.RegistryCronOnceHandler(...)       # 服务注入(注册一次性handler)
cron.RegistryCronCustomHandler(...)     # 服务注入(注册自定义handler)
```

# 示例
```gopackage main

import (
	"github.com/zly-app/service/cron"
	"github.com/zly-app/zapp"
)

func main() {
	// 注册cron服务
	cron.RegistryService()
	// 启用cron服务
	app := zapp.NewApp("test", cron.WithCronService())
	// 注册handler
	cron.RegistryCronHandler(app, "c1", "@every 1s", true, func(ctx *cron.Context) error {
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
