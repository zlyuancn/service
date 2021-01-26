/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package cron

import (
	"time"

	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
	"github.com/zlyuancn/zscheduler"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "cron"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 注册服务
func RegistryService(serviceType ...core.ServiceType) {
	t := DefaultServiceType
	if len(serviceType) > 0 {
		t = serviceType[0]
	}
	nowServiceType = t
	service.RegisterCreatorFunc(t, func(app core.IApp) core.IService {
		return NewCronService(app)
	})
}

// 启用cron服务
func WithCronService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册cron的Handler
func RegistryCronHandler(app core.IApp, name string, expression string, enable bool, handler RegistryCronHandlerFunc) {
	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  zscheduler.NewCronTrigger(expression),
		Executor: zscheduler.NewExecutor(0, 0, 1),
		Handler:  wrapHandler(handler),
		Enable:   enable,
	})
	app.InjectService(nowServiceType, task)
}

// 注册一次性cron的Handler
func RegistryCronOnceHandler(app core.IApp, name string, t time.Time, enable bool, handler RegistryCronHandlerFunc) {
	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  zscheduler.NewOnceTrigger(t),
		Executor: zscheduler.NewExecutor(0, 0, 1),
		Handler:  wrapHandler(handler),
		Enable:   enable,
	})
	app.InjectService(nowServiceType, task)
}

func RegistryCronCustomHandler(app core.IApp, name string, conf HandlerConfig) {
	var trigger zscheduler.ITrigger
	if conf.Expression != "" {
		trigger = zscheduler.NewCronTrigger(conf.Expression)
	} else {
		trigger = zscheduler.NewOnceTrigger(conf.OnceTime)
	}

	task := zscheduler.NewTaskOfConfig(name, zscheduler.TaskConfig{
		Trigger:  trigger,
		Executor: zscheduler.NewExecutor(conf.RetryCount, conf.RetryInterval, conf.MaxSyncExecuteCount),
		Handler:  wrapHandler(conf.Handler),
		Enable:   conf.Enable,
	})
	app.InjectService(nowServiceType, task)
}
