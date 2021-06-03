/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package middleware

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/core"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
)

func LoggerMiddleware(app core.IApp) iris.Handler {
	logResultInDevelop := &config.Conf.ShowApiResultInDevelop
	return func(ctx iris.Context) {
		log := app.NewMirrorLogger(ctx.Method(), ctx.Path())
		utils.Context.SaveLoggerToIrisContext(ctx, log)

		startTime := time.Now()
		log.Debug("api.request", zap.String("query", ctx.Request().URL.RawQuery))

		ctx.Next()

		latency := time.Since(startTime)
		fields := []interface{}{
			"api.response", zap.String("query", ctx.Request().URL.RawQuery),
			zap.String("latency_text", latency.String()),
			zap.Duration("latency", latency),
			zap.String("ip", ctx.RemoteAddr()),
		}

		if err, ok := ctx.Values().Get("error").(error); ok {
			if err == nil {
				err = fmt.Errorf("err{nil}")
			}
			fields = append(fields, zap.Error(err))
			log.Error(fields...)
		} else {
			if *logResultInDevelop {
				fields = append(fields, zap.Any("result", ctx.Values().Get("result")))
			}
			log.Debug(fields...)
		}
	}
}
