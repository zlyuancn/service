package middleware

import (
	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"github.com/opentracing/opentracing-go"
	"github.com/zly-app/zapp/core"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
)

// 用于构建相关log, trace等基础数据
func BaseMiddleware(app core.IApp, conf *config.Config) iris.Handler {
	return func(irisCtx *iris_context.Context) {
		name := irisCtx.Method() + ": " + irisCtx.Path()
		// 链路追踪
		span := opentracing.StartSpan(name)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(app.BaseContext(), span)
		utils.Context.SaveContextToIrisContext(irisCtx, ctx)

		// conf
		utils.Context.SaveConfToIrisContext(irisCtx, conf)

		// log
		log := app.NewTraceLogger(ctx)
		utils.Context.SaveLoggerToIrisContext(irisCtx, log)

		// handler
		irisCtx.Next()
	}
}
