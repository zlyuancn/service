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

// 注册服务
func RegistryService(serviceType ...core.ServiceType) {
	t := DefaultServiceType
	if len(serviceType) > 0 {
		t = serviceType[0]
	}
	nowServiceType = t
	service.RegisterCreatorFunc(t, func(app core.IApp) core.IService {
		return NewMysqlBinlogService(app)
	})
}

// 启用mysql-binlog服务
func WithMysqlBinlogService() zapp.Option {
	return zapp.WithService(nowServiceType)
}

// 注册handler
func RegistryMysqlBinlogHandler(app core.IApp, handler IEventHandler) {
	app.InjectService(nowServiceType, handler)
}
