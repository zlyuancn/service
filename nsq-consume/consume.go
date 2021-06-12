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
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/pkg/utils"
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"
)

type Context struct {
	core.ILogger
	*nsq.Message
	Topic               string
	Channel             string
	disableAutoRequeued bool // 关闭自动重排
}

// 关闭自动重排
func (ctx *Context) DisableAutoRequeued() {
	ctx.disableAutoRequeued = true
}

type RegistryNsqConsumerHandlerFunc = func(ctx *Context) error

type ConsumerConfig struct {
	Topic   string
	Channel string
	Handler RegistryNsqConsumerHandlerFunc
	Opts    []ConsumerOption
	*ServiceConfig
}

type consumerCli struct {
	app      core.IApp
	conf     *ConsumerConfig
	consumer *nsq.Consumer
	*consumerOptions
}

func newConsumer(app core.IApp, conf *ConsumerConfig) *consumerCli {
	c := &consumerCli{
		app:             app,
		conf:            conf,
		consumerOptions: newConsumerOptions(),
	}

	for _, o := range conf.Opts {
		o(c.consumerOptions)
	}

	if c.ConsumeAttempts == 0 {
		c.ConsumeAttempts = conf.ConsumeAttempts
	}
	return c
}

func (c *consumerCli) Start() error {
	if c.Disable {
		return nil
	}

	// 构建配置
	nsqConf := nsq.NewConfig()
	nsqConf.AuthSecret = c.conf.AuthSecret
	nsqConf.HeartbeatInterval = time.Duration(c.conf.HeartbeatInterval) * time.Millisecond
	nsqConf.ReadTimeout = time.Duration(c.conf.ReadTimeout) * time.Millisecond
	nsqConf.WriteTimeout = time.Duration(c.conf.WriteTimeout) * time.Millisecond
	nsqConf.DialTimeout = time.Duration(c.conf.DialTimeout) * time.Millisecond
	nsqConf.DefaultRequeueDelay = time.Duration(c.conf.RequeueDelay) * time.Millisecond
	nsqConf.MaxRequeueDelay = time.Duration(c.conf.MaxRequeueDelay) * time.Millisecond
	nsqConf.MaxInFlight = c.conf.MaxInFlight
	nsqConf.MaxAttempts = 0 // 解开sdk限制由我们实现

	// 创建消费者
	consumer, err := nsq.NewConsumer(c.conf.Topic, c.conf.Channel, nsqConf)
	if err != nil {
		return fmt.Errorf("创建nsq消费者失败, topic:%s, channel:%s, err:%s", c.conf.Topic, c.conf.Channel, err)
	}
	c.consumer = consumer

	// 添加消费handler
	c.consumer.AddConcurrentHandlers(c, utils.Ternary.Or(c.ThreadCount, c.conf.ThreadCount).(int))

	// 连接
	if c.conf.NsqLookupdAddress != "" {
		addresses := strings.Split(c.conf.NsqLookupdAddress, ",")
		return c.consumer.ConnectToNSQLookupds(addresses)
	}
	addresses := strings.Split(c.conf.NsqdAddress, ",")
	return c.consumer.ConnectToNSQDs(addresses)
}

func (c *consumerCli) Close() error {
	if c.Disable {
		return nil
	}

	c.consumer.Stop()
	<-c.consumer.StopChan
	return nil
}

func (c *consumerCli) HandleMessage(message *nsq.Message) error {
	ctx := &Context{
		ILogger: c.app.NewSessionLogger(
			zap.String("nsq_topic", c.conf.Topic),
			zap.String("nsq_channel", c.conf.Channel),
			zap.String("nsq_msg_id", string(message.ID[:])),
		),
		Message: message,
		Topic:   c.conf.Topic,
		Channel: c.conf.Channel,
	}

	ctx.Debug("nsqConsumer.receive")
	err := zutils.Recover.WrapCall(func() error {
		return c.conf.Handler(ctx)
	})

	if err == nil {
		ctx.Debug("nsqConsumer.success")
		return nil
	}

	// 如果关闭了自动重排
	if ctx.disableAutoRequeued {
		ctx.Error("nsqConsumer.error! and requeued is closed", zap.Error(err))
		return nil
	}

	// 检查自动重排次数
	if ctx.Attempts >= c.ConsumeAttempts {
		ctx.Error("nsqConsumer.error! reach the maximum automatic Requeue Attempts", zap.Error(err))
		return nil
	}

	ctx.Error("nsqConsumer.error!", zap.Error(err))
	return err
}
