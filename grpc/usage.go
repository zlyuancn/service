/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package grpc

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "grpc"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 启用grpc服务
func WithService(serviceType ...core.ServiceType) zapp.Option {
	if len(serviceType) > 0 && serviceType[0] != "" {
		nowServiceType = serviceType[0]
	}
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp, opts ...interface{}) core.IService {
		return NewGrpcService(app) // todo opts
	})
	return zapp.WithService(nowServiceType)
}

// 注册服务对象
func RegistryServerObject(app core.IApp, fn RegistryGrpcServiceFunc) {
	app.InjectService(nowServiceType, fn)
}
