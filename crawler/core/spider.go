package core

type ISpider interface {
	// 初始化
	Init() error
	// 提交初始化种子
	SubmitInitialSeed() []ISeed
	// 停止
	Stop() error
}
