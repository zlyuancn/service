/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package utils

import (
	"github.com/kataras/iris/v12"

	"github.com/zly-app/zapp/core"
)

var Context = new(contextUtil)

type contextUtil struct{}

const LoggerSaveFieldKey = "_api_logger"

// 将log保存在iris上下文中
func (c *contextUtil) SaveLoggerToIrisContext(ctx iris.Context, log core.ILogger) {
	ctx.Values().Set(LoggerSaveFieldKey, log)
}

// 从iris上下文中获取log, 如果失败会panic
func (c *contextUtil) MustGetLoggerFromIrisContext(ctx iris.Context) core.ILogger {
	return ctx.Values().Get(LoggerSaveFieldKey).(core.ILogger)
}
