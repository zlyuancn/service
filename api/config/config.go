/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package config

var Conf Config

const (
	// 默认bind
	defaultBind = ":8080"
	// 默认post允许最大数据大小
	defaultPostMaxMemory = 32 << 20
	// 将header中第一个有效的地址视为客户端地址
	defaultUseFirstValidRemoteAddrOfHeaders = true
)

// api服务配置
type Config struct {
	Bind                             string // bind地址
	IPWithNginxForwarded             bool   // 适配nginx的Forwarded获取ip, 优先级高于nginx的Real
	IPWithNginxReal                  bool   // 适配nginx的Real获取ip, 优先级高于sock连接的ip
	LogApiResultInDevelop            bool   // 在开发环境中输出api结果
	SendDetailedErrorInProduction    bool   // 在生产环境发送详细的错误到客户端
	AlwaysLogHeaders                 bool   // 总是输出headers日志, 如果设为false, 只会在出现错误时才会输出headers日志
	AlwaysLogBody                    bool   // 总是输出body日志, 如果设为false, 只会在出现错误时才会输出body日志
	UseFirstValidRemoteAddrOfHeaders bool   // 将header中第一个有效的地址视为客户端地址
	PostMaxMemory                    int64  // post允许客户端传输最大数据大小, 单位字节
}

func NewConfig() *Config {
	return &Config{
		UseFirstValidRemoteAddrOfHeaders: defaultUseFirstValidRemoteAddrOfHeaders,
	}
}

func (conf *Config) Check() {
	if conf.Bind == "" {
		conf.Bind = defaultBind
	}
	if conf.PostMaxMemory <= 0 {
		conf.PostMaxMemory = defaultPostMaxMemory
	}
}
