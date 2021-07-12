package api

import (
	"fmt"

	"github.com/zly-app/zapp/logger"
)

type options struct {
	Middlewares []interface{} // 中间件, 需要Wrap函数包装后才能用
}

type Option func(o *options)

func newOptions(opts ...interface{}) *options {
	o := &options{}
	for _, opt := range opts {
		fn, ok := opt.(Option)
		if !ok {
			logger.Log.Fatal(fmt.Sprintf("api 的选项必须传入 api.Option, 但收到了 %T", fn))
		}
		fn(o)
	}
	return o
}

// 添加全局中间件
func WithMiddleware(fn interface{}) Option {
	return func(o *options) {
		o.Middlewares = append(o.Middlewares, fn)
	}
}
