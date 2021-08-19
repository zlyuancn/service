package kafka_consume

type consumerOptions struct {
	Disable               bool // 禁用
	ConsumeCount          int
	ErrorsChannelCallback func(err error)
}

type ConsumerOption func(opts *consumerOptions)

func newConsumerOptions() *consumerOptions {
	return &consumerOptions{
		Disable:      false,
		ConsumeCount: 0,
	}
}

// 禁用
func WithConsumerDisable(disable ...bool) ConsumerOption {
	return func(opts *consumerOptions) {
		opts.Disable = len(disable) == 0 || disable[0]
	}
}

// 消费者数量, 建议设为topic的分区数, 0表示使用默认数量
func WithConsumerCount(consumeCount int) ConsumerOption {
	return func(opts *consumerOptions) {
		if consumeCount < 0 {
			consumeCount = 0
		}
		opts.ConsumeCount = consumeCount
	}
}

// 错误通道回调
func WithErrorsChannelCallback(fn func(err error)) ConsumerOption {
	return func(opts *consumerOptions) {
		opts.ErrorsChannelCallback = fn
	}
}
