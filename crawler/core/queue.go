package core

// 队列
type IQueue interface {
	/*
		**将种子放入队列
		 queueName 队列名
		 seed 种子
		 front 是否放在队列前面
	*/
	Put(queueName string, seed ISeed, front bool) error
	/*
		** 弹出一个种子
		 queueName 队列名
		 front 是否从队列前面弹出
	*/
	Pop(queueName string, front bool) (ISeed, error)
	// 检查队列是否为空
	CheckQueueIsEmpty() (bool, error)
	// 关闭
	Close() error
}
