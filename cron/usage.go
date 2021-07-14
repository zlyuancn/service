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

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t core.ServiceType) {
	nowServiceType = t
}

// 启用cron服务
func WithService() zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp) core.IService {
		return NewCronService(app)
	})
	return zapp.WithService(nowServiceType)
}

// 注册cron的Handler
func RegistryHandler(app core.IApp, name string, expression string, enable bool, handler Handler) {
	task := NewTaskOfConfig(name, TaskConfig{
		Trigger:  NewCronTrigger(expression),
		Executor: NewExecutor(0, 0, 1),
		Handler:  handler,
		Enable:   enable,
	})
	app.InjectService(nowServiceType, task)
}

// 注册一次性cron的Handler
func RegistryOnceHandler(app core.IApp, name string, t time.Time, enable bool, handler Handler) {
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
