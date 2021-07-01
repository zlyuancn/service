/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package mysql_binlog

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "mysql-binlog"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 启用mysql-binlog服务
func WithService(serviceType ...core.ServiceType) zapp.Option {
	if len(serviceType) > 0 && serviceType[0] != "" {
		nowServiceType = serviceType[0]
	}
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp, opts ...interface{}) core.IService {
		return NewMysqlBinlogService(app)
	})
	return zapp.WithService(nowServiceType)
}

// 注册handler
func RegistryHandler(app core.IApp, handler IEventHandler) {
	app.InjectService(nowServiceType, handler)
}
