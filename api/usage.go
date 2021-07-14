/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package api

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "api"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t core.ServiceType) {
	nowServiceType = t
}

// 启用app服务
func WithService(opts ...Option) zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp) core.IService {
		return NewHttpService(app, opts...)
	})
	return zapp.WithService(nowServiceType)
}

// 注册路由
func RegistryRouter(app core.IApp, fn ...RegisterApiRouterFunc) {
	a := make([]interface{}, len(fn))
	for i, h := range fn {
		a[i] = h
	}
	app.InjectService(nowServiceType, a...)
}
