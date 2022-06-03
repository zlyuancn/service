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
	// 是否启用请求数据校验
	defReqDataValidate = true
	// 是否对请求数据校验所有字段
	defReqDataValidateAllField = false
)

// grpc服务配置
type Config struct {
	Bind                    string // bind地址
	HeartbeatTime           int    // 心跳时间(毫秒), 默认20000
	EnableOpenTrace         bool   // 启用开放链路追踪
	ReqLogLevelIsInfo       bool   // 请求日志等级设为info
	RspLogLevelIsInfo       bool   // 响应日志等级设为info
	ReqDataValidate         bool   // 是否启用请求数据校验
	ReqDataValidateAllField bool   // 是否对请求数据校验所有字段. 如果设为true, 会对所有字段校验并返回所有的错误. 如果设为false, 校验错误会立即返回.
}

func newConfig() *Config {
	return &Config{
		EnableOpenTrace:         defaultEnableOpenTrace,
		ReqDataValidate:         defReqDataValidate,
		ReqDataValidateAllField: defReqDataValidateAllField,
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
