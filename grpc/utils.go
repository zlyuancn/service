/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"

	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
)

const loggerSaveFieldKey = "_grpc_logger"

// 基于传入的标准context生成一个新的标准context并保存log
func SaveLoggerToContext(ctx context.Context, log core.ILogger) context.Context {
	return context.WithValue(ctx, loggerSaveFieldKey, log)
}

// 从标准context中获取log
func GetLoggerFromContext(ctx context.Context) (core.ILogger, bool) {
	value := ctx.Value(loggerSaveFieldKey)
	log, ok := value.(core.ILogger)
	return log, ok
}

// 从标准context中获取log, 如果失败会panic
func MustGetLoggerFromContext(ctx context.Context) core.ILogger {
	log, ok := GetLoggerFromContext(ctx)
	if !ok {
		logger.Log.Panic("can't load logger from context")
	}
	return log
}
