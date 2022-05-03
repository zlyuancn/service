/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package grpc

const (
	// 默认bind
	defaultBind = ":3000"
	// 默认心跳时间
	defaultHeartbeatTime = 20000
	// 默认启用开放链路追踪
	defaultEnableOpenTrace = true
)

// grpc服务配置
type Config struct {
	Bind              string // bind地址
	HeartbeatTime     int    // 心跳时间(毫秒), 默认20000
	EnableOpenTrace   bool   // 启用开放链路追踪
	ReqLogLevelIsInfo bool   // 请求日志等级设为info
	RspLogLevelIsInfo bool   // 响应日志等级设为info
}

func newConfig() *Config {
	return &Config{
		EnableOpenTrace: defaultEnableOpenTrace,
	}
}

func (conf *Config) Check() {
	if conf.Bind == "" {
		conf.Bind = defaultBind
	}
	if conf.HeartbeatTime < 1000 {
		conf.HeartbeatTime = defaultHeartbeatTime
	}
}
