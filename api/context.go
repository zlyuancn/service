/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/30
   Description :
-------------------------------------------------
*/

package api

import (
	"context"
	"reflect"

	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/core"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
	"github.com/zly-app/service/api/validator"
)

type IrisContext = iris_context.Context

type Context struct {
	*IrisContext // 原始 iris.Context
	core.ILogger
	ctx  context.Context
	conf *config.Config
}

func makeContext(irisCtx iris.Context) *Context {
	return &Context{
		IrisContext: irisCtx,
		ILogger:     utils.Context.MustGetLoggerFromIrisContext(irisCtx),
		ctx:         utils.Context.MustGetContextFromIrisContext(irisCtx),
		conf:        utils.Context.MustGetConfFromIrisContext(irisCtx),
	}
}

//  bind api数据, 它会将api数据反序列化到a中, 如果a是结构体会验证a
func (c *Context) Bind(a interface{}) error {
	if err := c.ReadBody(a); err != nil {
		return ParamError.WithError(err)
	}

	if c.conf.BindLogLevelIsInfo {
		c.Info("api.request.bind", zap.Any("arg", a))
	} else {
		c.Debug("api.request.bind", zap.Any("arg", a))
	}

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
	return utils.Context.GetRemoteIP(c.IrisContext)
}

// 获取ctx, 这个ctx基于app.BaseContext并带链路追踪跨度
func (c *Context) Context() context.Context {
	return c.ctx
}
