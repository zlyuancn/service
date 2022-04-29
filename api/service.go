/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package api

import (
	"context"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/service/api/config"
	"github.com/zly-app/service/api/middleware"
)

type Party = iris.Party

// api注入函数定义
type RegisterApiRouterFunc = func(c core.IComponent, router Party)

type ApiService struct {
	app  core.IApp
	conf *config.Config
	*iris.Application
}

func NewApiService(app core.IApp, conf *config.Config, opts ...Option) *ApiService {
	// 处理选项
	o := newOptions(opts...)

	// irisApp
	irisApp := iris.New()
	irisApp.Logger().SetLevel("disable") // 关闭默认日志
	irisApp.Use(
		middleware.LoggerMiddleware(app, conf),
		cors.AllowAll(),
		middleware.Recover(),
	)
	irisApp.AllowMethods(iris.MethodOptions)

	// 配置项
	irisApp.Configure(o.Configurator...)

	// 中间件
	for _, fn := range o.Middlewares {
		irisApp.Use(WrapMiddleware(fn))
	}

	return &ApiService{
		app:         app,
		conf:        conf,
		Application: irisApp,
	}
}

func (a *ApiService) Start() error {
	a.app.Debug("正在启动api服务", zap.String("bind", a.conf.Bind))
	opts := []iris.Configurator{
		iris.WithoutBodyConsumptionOnUnmarshal,       // 重复消费
		iris.WithoutPathCorrection,                   // 不自动补全斜杠
		iris.WithOptimizations,                       // 启用性能优化
		iris.WithoutStartupLog,                       // 不要打印iris启动信息
		iris.WithPathEscape,                          // 解析path转义
		iris.WithFireMethodNotAllowed,                // 路由未找到时返回405而不是404
		iris.WithPostMaxMemory(a.conf.PostMaxMemory), // post允许客户端传输最大数据大小
	}
	if a.conf.IPWithNginxForwarded {
		opts = append(opts, iris.WithRemoteAddrHeader("X-Forwarded-For"))
	}
	if a.conf.IPWithNginxReal {
		opts = append(opts, iris.WithRemoteAddrHeader("X-Real-IP"))
	}
	return a.Run(iris.Addr(a.conf.Bind), opts...)
}

// 注册路由
func (a *ApiService) RegistryRouter(fn ...RegisterApiRouterFunc) {
	for _, h := range fn {
		h(a.app.GetComponent(), a.Party("/"))
	}
}

func (a *ApiService) Close() error {
	err := a.Shutdown(context.Background())
	a.app.Debug("api服务已关闭")
	return err
}
