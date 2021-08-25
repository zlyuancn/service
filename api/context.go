/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"net"
	"reflect"
	"strings"

	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/core"

	"github.com/zly-app/service/api/utils"
	"github.com/zly-app/service/api/validator"
)

type Context struct {
	*iris_context.Context // 原始 iris.Context
	core.ILogger
	OpenTraceSpan opentracing.Span // 链路追踪跨度
}

func makeContext(ctx iris.Context) *Context {
	return &Context{
		Context:       ctx,
		ILogger:       utils.Context.MustGetLoggerFromIrisContext(ctx),
		OpenTraceSpan: utils.Context.MustGetOpenTraceSpanFromIrisContext(ctx),
	}
}

//  bind api数据, 它会将api数据反序列化到a中, 如果a是结构体会验证a
func (c *Context) Bind(a interface{}) error {
	if err := c.ReadBody(a); err != nil {
		return ParamError.WithError(err)
	}

	c.Debug("api.request.bind", zap.Any("arg", a))

	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	err := validator.Valid(a)
	if err != nil {
		return ParamError.WithError(err)
	}
	return nil
}

// 试图解析并返回真实客户端的请求IP
func (c *Context) RemoteAddr() string {
	remoteHeaders := c.Application().ConfigurationReadOnly().GetRemoteAddrHeaders()
	for _, headerName := range remoteHeaders {
		ipAddresses := strings.Split(c.GetHeader(headerName), ",")
		for _, addr := range ipAddresses {
			if net.ParseIP(addr) != nil {
				return addr
			}
		}
	}

	addr := strings.TrimSpace(c.Request().RemoteAddr)
	if addr != "" {
		if ip, _, err := net.SplitHostPort(addr); err == nil {
			return ip
		}
	}
	return addr
}
