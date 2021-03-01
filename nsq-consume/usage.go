/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package nsq_consume

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "nsq-consume"

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
		return NewNsqConsumerService(app)
	})
}

// 启用nsq-consume服务
func WithNsqConsumeService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册handler
func RegistryNsqConsumeHandler(app core.IApp, topic, channel string, handler RegistryNsqConsumerHandlerFunc, opts ...HandlerOption) {
	app.InjectService(nowServiceType, newHandlerConfig(topic, channel, handler, opts...))
}
