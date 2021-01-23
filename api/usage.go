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

// 注册服务
func RegistryService(serviceType ...core.ServiceType) {
	t := DefaultServiceType
	if len(serviceType) > 0 {
		t = serviceType[0]
	}
	nowServiceType = t
	service.RegisterCreatorFunc(t, func(app core.IApp) core.IService {
		return NewHttpService(app)
	})
}

// 启用app服务
func WithApiService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册路由
func RegistryApiRouter(app core.IApp, fn RegisterApiRouterFunc) {
	app.InjectService(nowServiceType, fn)
}
