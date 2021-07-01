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
	"strings"

	"github.com/kataras/iris/v12"
	app_config "github.com/zly-app/zapp/config"
	app_utils "github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/utils"
)

func Recover() iris.Handler {
	isDebug := &app_config.Conf.Config().Frame.Debug
	sendDetailedErrorInProduction := &config.Conf.SendDetailedErrorInProduction
	isJson := &app_config.Conf.Config().Frame.Log.Json
	return func(ctx iris.Context) {
		err := app_utils.Recover.WrapCall(func() error {
			ctx.Next()
			return nil
		})
		if err == nil {
			return
		}

		if ctx.IsStopped() { // handled by other middleware.
			return
		}

		ctx.Values().Set("error", err)

		handlerName := ctx.Values().GetStringDefault("_handler_name", ctx.HandlerName())
		panicErrDetail := app_utils.Recover.GetRecoverErrorDetail(err)
		panicErrInfos := strings.Split(panicErrDetail, "\n")

		log := utils.Context.MustGetLoggerFromIrisContext(ctx)
		if *isJson {
			log.Error("panic",
				zap.String("handler_name", handlerName),
				zap.String("error", panicErrInfos[0]),
				zap.Strings("detail", panicErrInfos[1:]),
			)
		} else {
			var infoBuff bytes.Buffer
			infoBuff.WriteString("panic\n")
			infoBuff.WriteString("Recovered from a route's Handler: ")
			infoBuff.WriteString(handlerName)
			infoBuff.WriteByte('\n')
			infoBuff.WriteString(panicErrDetail)
			infoBuff.WriteByte('\n')
			log.Error(infoBuff.String())
		}

		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  "service internal error",
		}
		if *isDebug || *sendDetailedErrorInProduction {
			result["err_msg"] = append(
				[]string{fmt.Sprintf("Recovered from a route's Handler('%s')", handlerName)},
				panicErrInfos...,
			)
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}
