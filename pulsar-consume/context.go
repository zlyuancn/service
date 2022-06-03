package pulsar_consume

import (
	"context"

	"github.com/zly-app/zapp/core"
)

type ConsumerHandler func(ctx *Context) error

type Context struct {
	core.ILogger
	Ctx              context.Context // 上下文
	Msg              Message         // 消息
	SubscriptionName string          // 订阅名
}
