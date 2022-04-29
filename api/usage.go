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
	"go.uber.org/zap"

	"github.com/zly-app/service/api/config"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "api"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t core.ServiceType) {
	nowServiceType = t
}

// 启用app服务
func WithService(opts ...Option) zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp) core.IService {
		return newService(app, opts...)
	})
	return zapp.WithService(nowServiceType)
}

// 注册路由
func RegistryRouter(fn ...RegisterApiRouterFunc) {
	a := make([]interface{}, len(fn))
	for i, h := range fn {
		a[i] = h
	}
	zapp.App().InjectService(nowServiceType, a...)
}

type Service struct {
	app core.IApp
	api *ApiService
}

func (s *Service) Inject(a ...interface{}) {
	for _, h := range a {
		fn, ok := h.(RegisterApiRouterFunc)
		if !ok {
			s.app.Fatal("api服务注入类型错误, 它必须能转为 api.RegisterApiRouterFunc")
		}

		s.api.RegistryRouter(fn)
	}
}

func (s *Service) Start() error { return s.api.Start() }

func (s *Service) Close() error { return s.api.Close() }

func newService(app core.IApp, opts ...Option) *Service {
	conf := config.NewConfig()
	err := app.GetConfig().ParseServiceConfig(nowServiceType, conf, true)
	if err != nil {
		app.Fatal("获取api服务配置失败", zap.Error(err))
	}
	conf.Check()

	api := NewApiService(app, conf, opts...)
	return &Service{
		app: app,
		api: api,
	}
}
