package config

const (
	// 默认html编码
	defaultHtmlEncoding = "utf8"
)

type ServiceConfig struct {
	SpiderName   string // 爬虫名
	HtmlEncoding string // html编码
}

func NewConfig() *ServiceConfig {
	return &ServiceConfig{}
}
func (conf *ServiceConfig) Check() error {
	return nil
}
