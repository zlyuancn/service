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

// 注册服务
func RegistryService(serviceType ...core.ServiceType) {
	t := DefaultServiceType
	if len(serviceType) > 0 {
		t = serviceType[0]
	}
	nowServiceType = t
	service.RegisterCreatorFunc(t, func(app core.IApp, opts ...interface{}) core.IService {
		return NewGrpcService(app) // todo opts
	})
}

// 启用grpc服务
func WithGrpcService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册服务对象
func RegistryGrpcServerObject(app core.IApp, fn RegistryGrpcServiceFunc) {
	app.InjectService(nowServiceType, fn)
}
