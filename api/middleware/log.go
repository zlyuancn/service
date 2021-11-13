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
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"github.com/opentracing/opentracing-go"
	open_log "github.com/opentracing/opentracing-go/log"
	app_config "github.com/zly-app/zapp/config"
	app_utils "github.com/zly-app/zapp/pkg/utils"
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
	if app_config.Conf.Config().Frame.Log.Json {
		return loggerMiddlewareWithJson(app)
	}
	return loggerMiddleware(app)
}

// 以文本方式输出
func loggerMiddleware(app core.IApp) iris.Handler {
	isDebug := app_config.Conf.Config().Frame.Debug
	return func(ctx iris.Context) {
		startTime := time.Now()
		ip := utils.Context.GetRemoteIP(ctx)

		// log
		log := app.NewSessionLogger()
		utils.Context.SaveLoggerToIrisContext(ctx, log)

		// 链路追踪
		span := opentracing.StartSpan("api")
		defer span.Finish()
		c := opentracing.ContextWithSpan(app.BaseContext(), span)
		utils.Context.SaveContextToIrisContext(ctx, c)

		params := valuesToTexts(ctx.Request().URL.Query(), "=")

		// request
		span.SetTag("method", ctx.Method())
		span.SetTag("path", ctx.Path())
		span.SetTag("params", strings.Join(params, "\n"))
		span.SetTag("ip", ctx.RemoteAddr())
		var msgBuff bytes.Buffer
		msgBuff.WriteString("api.request path: ")
		msgBuff.WriteString(ctx.Method())
		msgBuff.WriteByte(' ')
		msgBuff.WriteString(ctx.Path())
		msgBuff.WriteString("\nparams:\n")
		for _, s := range params {
			msgBuff.WriteString("  ")
			msgBuff.WriteString(s)
			msgBuff.WriteByte('\n')
		}
		msgBuff.WriteByte('\n')
		log.Debug(msgBuff.String(), zap.String("ip", ip))

		// handler
		ctx.Next()

		// response
		msgBuff.Reset()
		msgBuff.WriteString("api.response path: ")
		msgBuff.WriteString(ctx.Method())
		msgBuff.WriteByte(' ')
		msgBuff.WriteString(ctx.Path())
		msgBuff.WriteString("\nparams:\n")
		for _, s := range params {
			msgBuff.WriteString("  ")
			msgBuff.WriteString(s)
			msgBuff.WriteByte('\n')
		}
		msgBuff.WriteByte('\n')

		latency := time.Since(startTime)
		fields := []interface{}{
			zap.String("ip", ip),
			zap.String("latency_text", latency.String()),
			zap.Duration("latency", latency),
		}

		// error
		err, hasErr := ctx.Values().Get("error").(error)
		hasPanic, _ := ctx.Values().Get("panic").(bool)
		if hasErr {
			if err == nil {
				err = fmt.Errorf("err{nil}")
			}
		}

		// headers
		if hasErr || config.Conf.AlwaysLogHeaders {
			headers := valuesToTexts(ctx.Request().Header, ": ")
			span.SetTag("headers", strings.Join(headers, "\n"))
			msgBuff.WriteString("headers:\n")
			for _, s := range headers {
				msgBuff.WriteString("  ")
				msgBuff.WriteString(s)
				msgBuff.WriteByte('\n')
			}
			msgBuff.WriteByte('\n')
		}

		// body
		if hasErr || config.Conf.AlwaysLogBody {
			var bodyText string
			if ctx.GetContentTypeRequested() == iris_context.ContentBinaryHeaderValue { // 流
				bodyText = fmt.Sprintf("body<bytesLen=%d>", ctx.GetContentLength())
			} else if ctx.GetContentLength() > config.Conf.LogBodyMaxSize { // 超长
				bodyText = fmt.Sprintf("body<len=%d>", ctx.GetContentLength())
			} else {
				body, _ := ctx.GetBody()
				bodyText = string(body)
			}
			span.LogFields(open_log.String("body", bodyText))
			msgBuff.WriteString("body:")
			msgBuff.WriteString(bodyText)
			msgBuff.WriteString("\n\n")
		}

		// result
		if !hasErr {
			var result string
			contentType := iris_context.TrimHeaderValue(ctx.ResponseWriter().Header().Get(iris_context.ContentTypeHeaderKey))
			if contentType == iris_context.ContentBinaryHeaderValue { // 流
				result = fmt.Sprintf("result<bytesLen=%d>", ctx.ResponseWriter().Written())
			} else if ctx.ResponseWriter().Written() > config.Conf.LogApiResultMaxSize { // 超长
				result = fmt.Sprintf("result<len=%d>", ctx.ResponseWriter().Written())
			} else {
				switch v := ctx.Values().Get("result").(type) {
				case nil:
					result = "result<nil>"
				case string:
					result = v
				default:
					result, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(v)
				}
			}
			span.LogFields(open_log.String("result", result))
			if (isDebug && config.Conf.LogApiResultInDevelop) || (!isDebug && config.Conf.LogApiResultInProd) {
				msgBuff.WriteString("result: ")
				msgBuff.WriteString(result)
				msgBuff.WriteString("\n\n")
			}
			log.Debug(append([]interface{}{msgBuff.String()}, fields...)...)
			return
		}

		// error
		if !hasPanic {
			span.SetTag("error", true)
			span.LogFields(open_log.Error(err))
			msgBuff.WriteString("err: ")
			msgBuff.WriteString(err.Error())
			msgBuff.WriteString("\n\n")
			log.Error(append([]interface{}{msgBuff.String()}, fields...)...)
			return
		}

		handlerName := ctx.Values().GetStringDefault("_handler_name", ctx.HandlerName())
		// panic
		panicErrDetail := app_utils.Recover.GetRecoverErrorDetail(err)
		panicErrInfos := strings.Split(panicErrDetail, "\n")
		span.SetTag("error", true)
		span.SetTag("panic", true)
		span.SetTag("handler_name", handlerName)
		span.SetTag("err", err.Error())
		span.LogFields(open_log.String("detail", panicErrDetail))

		msgBuff.WriteString("panic:\n")
		msgBuff.WriteString("  Recovered from a route's Handler: ")
		msgBuff.WriteString(handlerName)
		msgBuff.WriteString("  ")
		msgBuff.WriteString(strings.Join(panicErrInfos, "\n  "))
		msgBuff.WriteString("\n\n")
		log.Error(append([]interface{}{msgBuff.String()}, fields...)...)

		// send_error_result
		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  "service internal error",
		}
		if isDebug || config.Conf.SendDetailedErrorInProduction {
			result["err_msg"] = append(
				[]string{fmt.Sprintf("Recovered from a route's Handler: %s", handlerName)},
				panicErrInfos...,
			)
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}

// 以json方式输出
func loggerMiddlewareWithJson(app core.IApp) iris.Handler {
	isDebug := app_config.Conf.Config().Frame.Debug
	return func(ctx *iris_context.Context) {
		startTime := time.Now()
		ip := utils.Context.GetRemoteIP(ctx)

		// log
		log := app.NewSessionLogger()
		utils.Context.SaveLoggerToIrisContext(ctx, log)

		// 链路追踪
		span := opentracing.StartSpan("api")
		defer span.Finish()
		c := opentracing.ContextWithSpan(app.BaseContext(), span)
		utils.Context.SaveContextToIrisContext(ctx, c)

		params := valuesToTexts(ctx.Request().URL.Query(), "=")

		// request
		span.SetTag("method", ctx.Method())
		span.SetTag("path", ctx.Path())
		span.SetTag("params", strings.Join(params, "\n"))
		span.SetTag("ip", ctx.RemoteAddr())
		log.Debug(
			"api.request",
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.Strings("params", params),
			zap.String("ip", ip),
		)

		// handler
		ctx.Next()

		// response
		latency := time.Since(startTime)
		fields := []interface{}{
			"api.response",
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.Strings("params", params),
			zap.String("ip", ip),
			zap.String("latency_text", latency.String()),
			zap.Duration("latency", latency),
		}

		// error
		err, hasErr := ctx.Values().Get("error").(error)
		hasPanic, _ := ctx.Values().Get("panic").(bool)
		if hasErr {
			if err == nil {
				err = fmt.Errorf("err{nil}")
			}
		}

		// headers
		if hasErr || config.Conf.AlwaysLogHeaders {
			headers := valuesToTexts(ctx.Request().Header, ": ")
			span.SetTag("headers", strings.Join(headers, "\n"))
			fields = append(fields, zap.Strings("headers", headers))
		}

		// body
		if hasErr || config.Conf.AlwaysLogBody {
			var bodyText string
			if ctx.GetContentTypeRequested() == iris_context.ContentBinaryHeaderValue { // 流
				bodyText = fmt.Sprintf("body<bytesLen=%d>", ctx.GetContentLength())
			} else if ctx.GetContentLength() > config.Conf.LogBodyMaxSize { // 超长
				bodyText = fmt.Sprintf("body<len=%d>", ctx.GetContentLength())
			} else {
				body, _ := ctx.GetBody()
				bodyText = string(body)
			}
			span.LogFields(open_log.String("body", bodyText))
			fields = append(fields, zap.String("body", bodyText))
		}

		// result
		if !hasErr {
			var result string
			contentType := iris_context.TrimHeaderValue(ctx.ResponseWriter().Header().Get(iris_context.ContentTypeHeaderKey))
			if contentType == iris_context.ContentBinaryHeaderValue { // 流
				result = fmt.Sprintf("result<bytesLen=%d>", ctx.ResponseWriter().Written())
			} else if ctx.ResponseWriter().Written() > config.Conf.LogApiResultMaxSize { // 超长
				result = fmt.Sprintf("result<len=%d>", ctx.ResponseWriter().Written())
			} else {
				switch v := ctx.Values().Get("result").(type) {
				case nil:
					result = "result<nil>"
				case string:
					result = v
				default:
					result, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(v)
				}
			}
			span.LogFields(open_log.String("result", result))
			if (isDebug && config.Conf.LogApiResultInDevelop) || (!isDebug && config.Conf.LogApiResultInProd) {
				fields = append(fields, zap.String("result", result))
			}
			log.Debug(fields...)
			return
		}

		// error
		if !hasPanic {
			span.SetTag("error", true)
			span.SetTag("err", err.Error())
			fields = append(fields, zap.String("err", err.Error()))
			log.Error(fields...)
			return
		}

		handlerName := ctx.Values().GetStringDefault("_handler_name", ctx.HandlerName())
		// panic
		panicErrDetail := app_utils.Recover.GetRecoverErrorDetail(err)
		panicErrInfos := strings.Split(panicErrDetail, "\n")
		span.SetTag("error", true)
		span.SetTag("panic", true)
		span.SetTag("handler_name", handlerName)
		span.SetTag("err", err.Error())
		span.LogFields(open_log.String("detail", panicErrDetail))

		fields = append(fields,
			zap.Bool("panic", true),
			zap.String("handler_name", handlerName),
			zap.String("err", err.Error()),
			zap.Strings("detail", panicErrInfos),
		)
		log.Error(fields...)

		// send_error_result
		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  "service internal error",
		}
		if isDebug || config.Conf.SendDetailedErrorInProduction {
			result["err_msg"] = append(
				[]string{fmt.Sprintf("Recovered from a route's Handler: %s", handlerName)},
				panicErrInfos...,
			)
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}
