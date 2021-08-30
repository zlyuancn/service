package config

import (
	"errors"

	zapp_core "github.com/zly-app/zapp/core"
)

const (
	// 默认启用cookie
	defaultSpiderEnableCookie = false
	// 默认html编码
	defaultSpiderHtmlEncoding = "utf8"
	// 默认请求方法
	defaultSpiderRequestMethod = "get"
	// 默认提交初始化种子的时机
	defaultSpiderSubmitInitialSeedOpportunity = "start"
)

type SpiderConfig struct {
	Name          string // 爬虫名
	EnableCookie  bool   // 是否启用cookie
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
}

func newSpiderConfig(app zapp_core.IApp) SpiderConfig {
	return SpiderConfig{
		Name:                         app.Name(),
		EnableCookie:                 defaultSpiderEnableCookie,
		SubmitInitialSeedOpportunity: defaultSpiderSubmitInitialSeedOpportunity,
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
