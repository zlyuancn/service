/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/25
   Description :
-------------------------------------------------
*/

package nsq_consume

import (
	"errors"
	"runtime"
)

const (
	// 默认心跳间隔
	defaultHeartbeatInterval = 30000
	// 默认读取超时
	defaultReadTimeout = 30000
	// 默认写入超时
	defaultWriteTimeout = 5000
	// 默认连接超时
	defaultDialTimeout = 2000
	// MaxInFlight
	defaultMaxInFlight = 1024
	// 默认延时时间
	defaultRequeueDelay = 60000
	// 默认最大延时时间
	defaultMaxRequeueDelay = 600000
	// 默认消费尝试次数
	defaultConsumeAttempts = 5
)

type ServiceConfig struct {
	NsqdAddress       string // nsqd地址, localhost1:4150,localhost2:4150
	NsqLookupdAddress string // nsq发现服务地址, 优先级高于NsqdAddress, localhost1:4161,localhost2:4161
	AuthSecret        string // 验证秘钥
	HeartbeatInterval int64  // 心跳间隔(毫秒), 不能超过ReadTimeout
	ReadTimeout       int64  // 超时(毫秒)
	WriteTimeout      int64  // 超时(毫秒)
	DialTimeout       int64  // 超时(毫秒)
	MaxInFlight       int    // Maximum number of messages to allow in flight (concurrency knob)
	// 默认线程数, 默认为0表示使用逻辑cpu数量
	//
	// 同时处理信息的goroutine数
	ThreadCount     int
	RequeueDelay    int64  // 默认延时时间, 延时时间为-1时和消费失败自动发送延时消息时生效, 实际延时时间=延时时间x尝试次数(毫秒)
	MaxRequeueDelay int64  // 默认最大延时时间, 延时时间为-1时和消费失败自动发送延时消息时生效
	ConsumeAttempts uint16 // 消费尝试次数, 默认5, 最大65535
}

func newConfig() *ServiceConfig {
	return &ServiceConfig{
		NsqdAddress:       "",
		NsqLookupdAddress: "",
		AuthSecret:        "",
	}
}

func (conf *ServiceConfig) Check() error {
	if conf.ReadTimeout <= 0 {
		conf.ReadTimeout = defaultReadTimeout
	}
	if conf.WriteTimeout <= 0 {
		conf.WriteTimeout = defaultWriteTimeout
	}
	if conf.DialTimeout <= 0 {
		conf.DialTimeout = defaultDialTimeout
	}
	if conf.HeartbeatInterval <= 0 {
		conf.HeartbeatInterval = defaultHeartbeatInterval
	}
	if conf.HeartbeatInterval > conf.ReadTimeout {
		conf.HeartbeatInterval = conf.ReadTimeout
	}
	if conf.MaxInFlight <= 0 {
		conf.MaxInFlight = defaultMaxInFlight
	}
	if conf.ThreadCount <= 0 {
		conf.ThreadCount = runtime.NumCPU()
	}
	if conf.RequeueDelay <= 0 {
		conf.RequeueDelay = defaultRequeueDelay
	}
	if conf.MaxRequeueDelay <= 0 {
		conf.MaxRequeueDelay = defaultMaxRequeueDelay
	}
	if conf.ConsumeAttempts == 0 {
		conf.ConsumeAttempts = defaultConsumeAttempts
	}
	if conf.NsqdAddress == "" && conf.NsqLookupdAddress == "" {
		return errors.New("address为空")
	}

	return nil
}
