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
	service.RegisterCreatorFunc(t, func(app core.IApp, opts ...interface{}) core.IService {
		return NewCronService(app) // todo opts
	})
}

// 启用cron服务
func WithCronService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册cron的Handler
func RegistryCronHandler(app core.IApp, name string, expression string, enable bool, handler Handler) {
	task := NewTaskOfConfig(name, TaskConfig{
		Trigger:  NewCronTrigger(expression),
		Executor: NewExecutor(0, 0, 1),
		Handler:  handler,
		Enable:   enable,
	})
	app.InjectService(nowServiceType, task)
}

// 注册一次性cron的Handler
func RegistryCronOnceHandler(app core.IApp, name string, t time.Time, enable bool, handler Handler) {
	task := NewTaskOfConfig(name, TaskConfig{
		Trigger:  NewOnceTrigger(t),
		Executor: NewExecutor(0, 0, 1),
		Handler:  handler,
		Enable:   enable,
	})
	app.InjectService(nowServiceType, task)
}

// 注册自定义task
func RegistryTask(app core.IApp, task ITask) {
	app.InjectService(nowServiceType, task)
}
