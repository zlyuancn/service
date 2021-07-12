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

// 启用app服务
func WithService(serviceType ...core.ServiceType) zapp.Option {
	if len(serviceType) > 0 && serviceType[0] != "" {
		nowServiceType = serviceType[0]
	}
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp, opts ...interface{}) core.IService {
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
