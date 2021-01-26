/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/26
   Description :
-------------------------------------------------
*/

package cron

import (
	"fmt"

	"github.com/zly-app/zapp/core"
	"github.com/zlyuancn/zscheduler"
	"go.uber.org/zap"
)

type CronService struct {
	app       core.IApp
	scheduler zscheduler.IScheduler
}

func NewCronService(app core.IApp) core.IService {
	conf := newConfig()
	vi := app.GetConfig().GetViper()
	confKey := "services." + string(nowServiceType)
	if vi.IsSet(confKey) {
		if err := vi.UnmarshalKey(confKey, conf); err != nil {
			app.Fatal(fmt.Errorf("无法解析<%s>服务配置: %s", nowServiceType, err))
		}
	}
	conf.check()

	return &CronService{
		app: app,
		scheduler: zscheduler.NewScheduler(
			zscheduler.WithLogger(nil),
			zscheduler.WithGoroutinePool(conf.ThreadCount, conf.MaxTaskQueueSize),
			zscheduler.WithObserver(newObserver(app)),
		),
	}
}

func (c *CronService) Inject(a ...interface{}) {
	for _, v := range a {
		task, ok := v.(zscheduler.ITask)
		if !ok {
			c.app.Fatal("Cron服务注入类型错误, 它必须能转为 zscheduler.ITask")
		}

		if ok := c.scheduler.AddTask(task); !ok {
			c.app.Fatal("添加Cron任务失败, 可能是名称重复", zap.String("name", task.Name()))
		}
	}
}

func (c *CronService) Start() error {
	c.scheduler.Start()
	return nil
}

func (c *CronService) Close() error {
	c.scheduler.Stop()
	return nil
}
