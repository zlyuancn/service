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

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t core.ServiceType) {
	nowServiceType = t
}

// 启用grpc服务
func WithService() zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp) core.IService {
		return NewGrpcService(app)
	})
	return zapp.WithService(nowServiceType)
}

// 注册服务对象
func RegistryServerObject(fn RegistryGrpcServiceFunc) {
	zapp.App().InjectService(nowServiceType, fn)
}
