package api

import (
	"github.com/kataras/iris/v12"
)

type options struct {
	Middlewares  []interface{}       // 中间件, 需要Wrap函数包装后才能用
	Configurator []iris.Configurator // 配置项
}

type Option func(o *options)

func newOptions(opts ...Option) *options {
	o := &options{}
	for _, fn := range opts {
		fn(o)
	}
	return o
}

// 添加iris配置项
func WithConfigurator(configs ...iris.Configurator) Option {
	return func(o *options) {
		o.Configurator = append(o.Configurator, configs...)
	}
}

// 添加全局中间件
func WithMiddleware(fn interface{}) Option {
	return func(o *options) {
		o.Middlewares = append(o.Middlewares, fn)
	}
}
