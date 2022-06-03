/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/29
   Description :
-------------------------------------------------
*/

package middleware

import (
	"github.com/kataras/iris/v12"
	app_utils "github.com/zly-app/zapp/pkg/utils"
)

func Recover() iris.Handler {
	return func(ctx iris.Context) {
		err := app_utils.Recover.WrapCall(func() error {
			ctx.Next()
			return nil
		})
		if err != nil {
			ctx.Values().Set("error", err)
			ctx.Values().Set("panic", true)
		}
	}
}
