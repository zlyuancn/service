/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/1
   Description :
-------------------------------------------------
*/

package nsq_consume

type handlerOptions struct {
	ThreadCount int
}

type HandlerOption func(opts *handlerOptions)

func newHandlerOptions() *handlerOptions {
	return &handlerOptions{
		ThreadCount: 0,
	}
}

// 线程数, 默认为0表示使用配置的默认线程数
//
// 同时处理信息的goroutine数
func WithHandlerThreadCount(threadCount int) HandlerOption {
	return func(opts *handlerOptions) {
		if threadCount < 0 {
			threadCount = 0
		}
		opts.ThreadCount = threadCount
	}
}
