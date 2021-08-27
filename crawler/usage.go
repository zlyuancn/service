package crawler

import (
	"github.com/zly-app/zapp"
	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/service"

	"github.com/zly-app/service/crawler/core"
)

// 默认服务类型
const DefaultServiceType zapp_core.ServiceType = "crawler"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t zapp_core.ServiceType) {
	nowServiceType = t
}

// 启用crawler服务
func WithService() zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app zapp_core.IApp) zapp_core.IService {
		return NewCrawlerService(app)
	})
	return zapp.WithService(nowServiceType)
}

// 注册spider
func RegistryHandler(spider core.ISpider) {
	zapp.App().InjectService(nowServiceType, spider)
}
