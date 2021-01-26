/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package cron

import (
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"github.com/zlyuancn/zscheduler"
)

const loggerSaveFieldKey = "_grpc_logger"

// 将log存入job, 如果meta不是nil或map[string]interface{}会panic
func SaveLoggerToJob(job zscheduler.IJob, log core.ILogger) {
	if job.Meta() == nil {
		job.SetMeta(map[string]interface{}{
			loggerSaveFieldKey: log,
		})
		return
	}
	job.Meta().(map[string]interface{})[loggerSaveFieldKey] = log
}

// 从job中获取log
func GetLoggerFromJob(job zscheduler.IJob) (core.ILogger, bool) {
	if job.Meta() == nil {
		return nil, false
	}
	log, ok := job.Meta().(map[string]interface{})[loggerSaveFieldKey]
	if !ok {
		return nil, false
	}
	l, ok := log.(core.ILogger)
	return l, ok
}

// 从job中获取log, 如果失败会panic
func MustGetLoggerFromJob(job zscheduler.IJob) core.ILogger {
	log, ok := GetLoggerFromJob(job)
	if !ok {
		logger.Log.Panic("can't load logger from job")
	}
	return log
}
