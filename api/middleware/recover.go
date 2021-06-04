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
	showDetailedErrorInProduction := &config.Conf.ShowDetailedErrorInProduction
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
		panicErrInfos := strings.Split(app_utils.Recover.GetRecoverErrorDetail(err), "\n")

		log := utils.Context.MustGetLoggerFromIrisContext(ctx)
		log.Error(strings.Join(panicErrInfos, "\n"), zap.Bool("panic", true), zap.String("handler_name", handlerName))

		result := map[string]interface{}{
			"err_code": 1,
			"err_msg":  "service internal error",
		}
		if *isDebug || *showDetailedErrorInProduction {
			result["err_msg"] = append(
				[]string{fmt.Sprintf("Recovered from a route's Handler('%s')", handlerName)},
				panicErrInfos...,
			)
		}
		_, _ = ctx.JSON(result)
		ctx.StopExecution()
	}
}
