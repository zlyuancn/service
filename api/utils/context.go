/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package utils

import (
	"net"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/opentracing/opentracing-go"

	"github.com/zly-app/zapp/core"
)

var Context = new(contextUtil)

type contextUtil struct{}

// 日志保存字段
const LoggerSaveFieldKey = "_api_logger"

// 链路追踪跨度保存字段
const OpenTraceSpanSaveFieldKey = "_open_trace_span"

// 将log保存在iris上下文中
func (c *contextUtil) SaveLoggerToIrisContext(ctx iris.Context, log core.ILogger) {
	ctx.Values().Set(LoggerSaveFieldKey, log)
}

// 从iris上下文中获取log, 如果失败会panic
func (c *contextUtil) MustGetLoggerFromIrisContext(ctx iris.Context) core.ILogger {
	return ctx.Values().Get(LoggerSaveFieldKey).(core.ILogger)
}

// 将链路追踪跨度保存在iris上下文中
func (c *contextUtil) SaveOpenTraceSpanToIrisContext(ctx iris.Context, span opentracing.Span) {
	ctx.Values().Set(OpenTraceSpanSaveFieldKey, span)
}

// 从iris上下文中获取链路追踪跨度, 如果失败会panic
func (c *contextUtil) MustGetOpenTraceSpanFromIrisContext(ctx iris.Context) opentracing.Span {
	return ctx.Values().Get(OpenTraceSpanSaveFieldKey).(opentracing.Span)
}

// 试图解析并返回真实客户端的请求IP
func (c *contextUtil) GetRemoteIP(ctx iris.Context) string {
	remoteHeaders := ctx.Application().ConfigurationReadOnly().GetRemoteAddrHeaders()
	for _, headerName := range remoteHeaders {
		ipAddresses := strings.Split(ctx.GetHeader(headerName), ",")
		for _, addr := range ipAddresses {
			if net.ParseIP(addr) != nil {
				return addr
			}
		}
	}

	addr := strings.TrimSpace(ctx.Request().RemoteAddr)
	if addr != "" {
		if ip, _, err := net.SplitHostPort(addr); err == nil {
			return ip
		}
	}
	return addr
}
