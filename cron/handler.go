/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/26
   Description :
-------------------------------------------------
*/

package cron

import (
	"time"

	"github.com/zly-app/zapp/core"
	"github.com/zlyuancn/zscheduler"
)

type RegistryCronHandlerFunc = func(ctx *Context) error

type HandlerConfig struct {
	// cron表达式, 优先级高于OnceTime
	Expression string
	// 一次性触发时间
	OnceTime time.Time
	// 失败时重试次数
	RetryCount int64
	// 失败重试时间隔
	RetryInterval time.Duration
	// 最大同步执行数, 0表示不限
	//
	// 假设上一次任务还没执行完毕就到了下一次触发时间, 这时就存在2个任务同时执行, 这个参数表示这个任务能同时执行的个数
	MaxSyncExecuteCount int64
	// handler
	Handler RegistryCronHandlerFunc
	// 是否启用
	Enable bool
}

type Context struct {
	core.ILogger
	zscheduler.IJob
}

// 包装 handler 以适应 zscheduler.Handler
func wrapHandler(handler RegistryCronHandlerFunc) zscheduler.Handler {
	return func(job zscheduler.IJob) (err error) {
		ctx := &Context{
			ILogger: MustGetLoggerFromJob(job),
			IJob:    job,
		}
		return handler(ctx)
	}
}
