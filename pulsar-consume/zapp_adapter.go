package pulsar_consume

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/service"
	"go.uber.org/zap"
)

// 默认服务类型
const DefaultServiceType core.ServiceType = "pulsar-consume"

// 当前服务类型
var nowServiceType = DefaultServiceType

// 设置服务类型, 这个函数应该在 zapp.NewApp 之前调用
func SetServiceType(t core.ServiceType) {
	nowServiceType = t
}

// 启用pulsar-consume服务
func WithService() zapp.Option {
	service.RegisterCreatorFunc(nowServiceType, func(app core.IApp) core.IService {
		return NewServiceAdapter(app)
	})
	return zapp.WithService(nowServiceType)
}

type ServiceAdapter struct {
	service *PulsarConsumeService
}

func (s *ServiceAdapter) Inject(a ...interface{}) {

}

func (s *ServiceAdapter) Start() error {
	s.service.Start()
	return nil
}

func (s *ServiceAdapter) Close() error {
	s.service.Close()
	return nil
}

func NewServiceAdapter(app core.IApp) core.IService {
	conf := NewConfig()
	err := app.GetConfig().ParseServiceConfig(nowServiceType, conf)
	if err != nil {
		logger.Log.Panic("服务配置错误", zap.String("serviceType", string(nowServiceType)), zap.Error(err))
	}
	s, err := NewConsumeService(app, conf)
	if err != nil {
		logger.Log.Panic("创建服务失败", zap.String("serviceType", string(nowServiceType)), zap.Error(err))
	}

	return &ServiceAdapter{s}
}
