/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package middleware

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/kataras/iris/v12"
	app_config "github.com/zly-app/zapp/config"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/core"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
)

func valuesToTexts(values map[string][]string, sep string) []string {
	var texts []string
	for k, vs := range values {
		for _, v := range vs {
			texts = append(texts, k+sep+v)
		}
	}
	sort.Strings(texts)
	return texts
}

func LoggerMiddleware(app core.IApp) iris.Handler {
	isDebug := &app_config.Conf.Config().Frame.Debug
	logResultInDevelop := &config.Conf.LogApiResultInDevelop
	isJson := &app_config.Conf.Config().Frame.Log.Json
	return func(ctx iris.Context) {
		startTime := time.Now()
		addr := ctx.RemoteAddr()

		log := app.NewSessionLogger(zap.String("method", ctx.Method()), zap.String("path", ctx.Path()))
		utils.Context.SaveLoggerToIrisContext(ctx, log)

		body, _ := ctx.GetBody()

		if *isJson {
			log.Debug(
				"api.request",
				zap.String("ip", addr),
				zap.Strings("headers", valuesToTexts(ctx.Request().Header, ": ")),
				zap.Strings("params", valuesToTexts(ctx.Request().URL.Query(), "=")),
				zap.String("body", string(body)),
			)
		} else {
			var infoBuff bytes.Buffer
			infoBuff.WriteString("api.request\nheaders:\n")
			for _, s := range valuesToTexts(ctx.Request().Header, ": ") {
				infoBuff.WriteString("  ")
				infoBuff.WriteString(s)
				infoBuff.WriteByte('\n')
			}
			infoBuff.WriteString("\nparams:\n")
			for _, s := range valuesToTexts(ctx.Request().URL.Query(), "=") {
				infoBuff.WriteString("  ")
				infoBuff.WriteString(s)
				infoBuff.WriteByte('\n')
			}
			infoBuff.WriteString("\nbody:")
			infoBuff.Write(body)
			infoBuff.WriteByte('\n')
			log.Debug(infoBuff.String(), zap.String("ip", addr))
		}

		ctx.Next()

		latency := time.Since(startTime)
		fields := []interface{}{
			"api.response",
			zap.String("query", ctx.Request().URL.RawQuery),
			zap.String("ip", addr),
			zap.String("latency_text", latency.String()),
			zap.Duration("latency", latency),
		}

		if err, ok := ctx.Values().Get("error").(error); ok {
			if err == nil {
				err = fmt.Errorf("err{nil}")
			}
			fields = append(fields, zap.Error(err))
			log.Error(fields...)
		} else {
			if *isDebug && *logResultInDevelop {
				fields = append(fields, zap.Any("result", ctx.Values().Get("result")))
			}
			log.Debug(fields...)
		}
	}
}
