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
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	app_config "github.com/zly-app/zapp/config"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
)

func Recover() iris.Handler {
	isDebug := &app_config.Conf.Config().Frame.Debug
	showDetailedErrorInProduction := &config.Conf.ShowDetailedErrorInProduction
	return func(ctx iris.Context) {
		err := WrapCall(func() error {
			ctx.Next()
			return nil
		})
		if err == nil {
			return
		}

		if ctx.IsStopped() { // handled by other middleware.
			return
		}

		var callers []string
		if re, ok := err.(RecoverError); ok {
			callers = make([]string, len(re.Callers()))
			for i, c := range re.Callers() {
				callers[i] = fmt.Sprintf("%s:%d", c.File, c.Line)
			}
		} else {
			callers = append(callers, err.Error())
		}

		handlerName := ctx.Values().GetStringDefault("_handler_name", ctx.HandlerName())
		logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", handlerName)
		logMessage += fmt.Sprint(getRequestLogs(ctx))
		logMessage += fmt.Sprintf("err: %s\n", err)
		logMessage += strings.Join(callers, "\n")
		log := utils.Context.MustGetLoggerFromIrisContext(ctx)
		log.Error(logMessage)
		ctx.Values().Set("error", err)

		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  strings.Split(logMessage, "\n"),
		}
		if !*isDebug && !*showDetailedErrorInProduction {
			result["err_msg"] = "service internal error"
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}

func getRequestLogs(ctx iris.Context) string {
	var status, ip, method, path string
	status = strconv.Itoa(ctx.GetStatusCode())
	path = ctx.Path()
	method = ctx.Method()
	ip = ctx.RemoteAddr()
	return fmt.Sprintf("%v %s %s %s\n", status, path, method, ip)
}
