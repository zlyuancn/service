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

// 启用nsq-consume服务
func WithService(serviceType ...core.ServiceType) zapp.Option {
	if len(serviceType) > 0 && serviceType[0] != "" {
		nowServiceType = serviceType[0]
	}
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp, opts ...interface{}) core.IService {
		return NewNsqConsumeService(app) // todo opts
	})
	return zapp.WithService(nowServiceType)
}

// 注册handler
func RegistryHandler(app core.IApp, topic, channel string, handler RegistryNsqConsumerHandlerFunc, opts ...ConsumerOption) {
	app.InjectService(nowServiceType, &ConsumerConfig{
		Topic:   topic,
		Channel: channel,
		Handler: handler,
		Opts:    opts,
	})
}
