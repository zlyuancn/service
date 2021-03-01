/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/25
   Description :
-------------------------------------------------
*/

package nsq_consume

import (
	"github.com/nsqio/go-nsq"
	"github.com/zly-app/zapp/core"
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"
)

type Context struct {
	core.ILogger
	*nsq.Message
	Topic   string
	Channel string
}
type RegistryNsqConsumerHandlerFunc = func(ctx *Context) error

type handlerConfig struct {
	app      core.IApp
	consumer *nsq.Consumer
	Topic    string
	Channel  string
	Handler  RegistryNsqConsumerHandlerFunc
	*handlerOptions
}

func newHandlerConfig(topic, channel string, handler RegistryNsqConsumerHandlerFunc, opts ...HandlerOption) *handlerConfig {
	hc := &handlerConfig{
		consumer:       nil,
		Topic:          topic,
		Channel:        channel,
		Handler:        handler,
		handlerOptions: newHandlerOptions(),
	}
	for _, o := range opts {
		o(hc.handlerOptions)
	}
	return hc
}

func (h *handlerConfig) SetConsumer(app core.IApp, consumer *nsq.Consumer) {
	h.app = app
	h.consumer = consumer
}

func (h *handlerConfig) HandleMessage(message *nsq.Message) error {
	ctx := &Context{
		ILogger: h.app.NewMirrorLogger(h.Topic, h.Channel, string(message.ID[:])),
		Message: message,
		Topic:   h.Topic,
		Channel: h.Channel,
	}

	ctx.Debug("nsqConsumer.receive")
	err := zutils.Recover.WrapCall(func() error {
		return h.Handler(ctx)
	})

	if err != nil {
		ctx.Error("nsqConsumer.error!", zap.Error(err))
	} else {
		ctx.Debug("nsqConsumer.success")
	}

	return err
}
