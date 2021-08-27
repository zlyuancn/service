package crawler

import (
	"fmt"

	zapp_core "github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zly-app/service/crawler/config"
	"github.com/zly-app/service/crawler/core"
)

type CrawlerService struct {
	app    zapp_core.IApp
	conf   *config.ServiceConfig
	spider core.ISpider
}

func (c *CrawlerService) Inject(a ...interface{}) {
	if c.spider != nil {
		c.app.Fatal("crawler服务重复注入")
	}

	if len(a) != 1 {
		c.app.Fatal("crawler服务注入数量必须为1个")
	}

	var ok bool
	c.spider, ok = a[0].(core.ISpider)
	if !ok {
		c.app.Fatal("crawler服务注入类型错误, 它必须能转为 crawler/core.ISpider")
	}
}

func (c *CrawlerService) Start() error {
	err := c.spider.Init()
	if err != nil {
		return fmt.Errorf("spider初始化失败: %v", err)
	}
	return nil
}

func (c *CrawlerService) Close() error {
	err := c.spider.Stop()
	if err != nil {
		c.app.Error("spider停止时出错: %v", err)
	}

	// todo 在这里关闭其他组件
	return nil
}

func NewCrawlerService(app zapp_core.IApp) zapp_core.IService {
	conf := config.NewConfig()
	err := app.GetConfig().ParseServiceConfig(nowServiceType, conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		logger.Log.Panic("服务配置错误", zap.String("serviceType", string(nowServiceType)), zap.Error(err))
	}

	// todo 这里初始化相关组件

	return &CrawlerService{
		app:  app,
		conf: conf,
	}
}
