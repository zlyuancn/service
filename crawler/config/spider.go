package config

import (
	"errors"

	zapp_core "github.com/zly-app/zapp/core"
)

const (
	// 默认启用cookie
	defaultSpiderUseCookie = false
	// 默认html编码
	defaultSpiderHtmlEncoding = "utf8"
	// 默认请求方法
	defaultSpiderRequestMethod = "get"
	// 默认提交初始化种子的时机
	defaultSpiderSubmitInitialSeedOpportunity = "start"
	// 默认使用调度器
	defaultSpiderUseScheduler = false
)

type SpiderConfig struct {
	Name          string // 爬虫名
	UseCookie     bool   // 是否启用cookie
	HtmlEncoding  string // 默认html编码
	RequestMethod string // 默认请求方法
	/*
		**提交初始化种子的时机
		 none 无
		 start 启动时
		 YYYY-MM-DD hh:mm:ss 指定时间触发
		 cron表达式
	*/
	SubmitInitialSeedOpportunity string
	// 使用调度器, 提交初始化种子的时机交给调度器管理, 这可以解决多进程运行时每个进程都在提交种子
	UseScheduler bool

	ExpectHttpStatusCode  []int // 期望的http状态码列表
	InvalidHttpStatusCode []int // 无效的http状态码列表, 如果配置了ExpectHttpStatusCode, 则以ExpectHttpStatusCode为准
}

func newSpiderConfig(app zapp_core.IApp) SpiderConfig {
	return SpiderConfig{
		Name:                         app.Name(),
		UseCookie:                    defaultSpiderUseCookie,
		SubmitInitialSeedOpportunity: defaultSpiderSubmitInitialSeedOpportunity,
		UseScheduler:                 defaultSpiderUseScheduler,
	}
}
func (conf *SpiderConfig) Check() error {
	if conf.Name == "" {
		return errors.New("spider.name is empty")
	}
	if conf.HtmlEncoding == "" {
		conf.HtmlEncoding = defaultSpiderHtmlEncoding
	}
	if conf.RequestMethod == "" {
		conf.RequestMethod = defaultSpiderRequestMethod
	}
	return nil
}
