/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/25
   Description :
-------------------------------------------------
*/

package nsq_consume

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/zly-app/zapp/core"
)

type NsqConsumeService struct {
	app      core.IApp
	handlers []*handlerConfig
}

func (n *NsqConsumeService) Inject(a ...interface{}) {
	for _, v := range a {
		h, ok := v.(*handlerConfig)
		if !ok {
			n.app.Fatal("nsq消费服务注入类型错误, 它必须能转为 *nsq_consume.handlerConfig")
		}
		n.handlers = append(n.handlers, h)
	}
}

func (n *NsqConsumeService) Start() error {
	// 加载配置
	conf := newConfig()
	err := n.app.GetConfig().ParseServiceConfig(nowServiceType, conf)
	if err != nil {
		return err
	}
	conf.Check()

	// 检查配置
	var nsqdAddress, nsqLookupdAddress []string
	if conf.NsqLookupdAddress != "" {
		nsqLookupdAddress = strings.Split(conf.NsqLookupdAddress, ",")
	} else if conf.NsqdAddress != "" {
		nsqdAddress = strings.Split(conf.NsqdAddress, ",")
	} else {
		return fmt.Errorf("nsq消费服务address为空")
	}

	// 创建消费者
	for _, h := range n.handlers {
		nsqConf := nsq.NewConfig()
		nsqConf.AuthSecret = conf.AuthSecret
		nsqConf.HeartbeatInterval = time.Duration(conf.HeartbeatInterval) * time.Millisecond
		nsqConf.ReadTimeout = time.Duration(conf.ReadTimeout) * time.Millisecond
		nsqConf.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Millisecond
		nsqConf.DialTimeout = time.Duration(conf.DialTimeout) * time.Millisecond

		consumer, err := nsq.NewConsumer(h.Topic, h.Channel, nsqConf)
		if err != nil {
			return fmt.Errorf("创建nsq消费者失败, topic:%s, channel:%s, err:%s", h.Topic, h.Channel, err)
		}
		h.SetConsumer(consumer)
		consumer.AddConcurrentHandlers(h, conf.ThreadCount)
	}

	// 连接
	for _, h := range n.handlers {
		if len(nsqLookupdAddress) > 0 {
			err = h.consumer.ConnectToNSQLookupds(nsqLookupdAddress)
		} else {
			err = h.consumer.ConnectToNSQDs(nsqdAddress)
		}
		if err != nil {
			return fmt.Errorf("nsq消费服务启动失败: %s", err)
		}
	}
	return nil
}

func (n *NsqConsumeService) Close() error {
	var wg sync.WaitGroup
	wg.Add(len(n.handlers))
	for _, h := range n.handlers {
		go func(consumer *nsq.Consumer) {
			defer wg.Done()
			consumer.Stop()
			<-consumer.StopChan
		}(h.consumer)
	}
	wg.Wait()
	return nil
}

func NewNsqConsumerService(app core.IApp) core.IService {
	return &NsqConsumeService{
		app: app,
	}
}
