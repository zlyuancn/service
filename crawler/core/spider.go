package core

type ISpider interface {
	// 初始化
	Init() error
	// 提交初始化种子
	SubmitInitialSeed(queue IQueue, queueName string, front bool) (int, error)
	// 停止
	Stop() error
}
