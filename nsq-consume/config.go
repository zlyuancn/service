/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/25
   Description :
-------------------------------------------------
*/

package nsq_consume

import (
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
	// 默认线程数
	defaultThreadCount = 0
)

type Config struct {
	NsqdAddress       string // nsqd地址, localhost1:4150,localhost2:4150
	NsqLookupdAddress string // nsq发现服务地址, 优先级高于NsqdAddress, localhost1:4161,localhost2:4161
	AuthSecret        string // 验证秘钥
	HeartbeatInterval int64  // 心跳间隔(毫秒), 不能超过ReadTimeout, 0表示无
	ReadTimeout       int64  // 超时(毫秒)
	WriteTimeout      int64  // 超时(毫秒)
	DialTimeout       int64  // 超时(毫秒)
	MaxInFlight       int    // Maximum number of messages to allow in flight (concurrency knob)
	// 线程数, 默认为0表示使用逻辑cpu数量
	//
	// 同时处理信息的goroutine数
	ThreadCount int
}

func newConfig() *Config {
	return &Config{
		NsqdAddress:       "",
		NsqLookupdAddress: "",
		AuthSecret:        "",
		HeartbeatInterval: defaultHeartbeatInterval,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		DialTimeout:       defaultDialTimeout,
		MaxInFlight:       defaultMaxInFlight,
		ThreadCount:       defaultThreadCount,
	}
}

func (conf *Config) Check() {
	if conf.HeartbeatInterval <= 0 {
		conf.HeartbeatInterval = -1
	}
	if conf.ReadTimeout <= 0 {
		conf.ReadTimeout = defaultReadTimeout
	}
	if conf.WriteTimeout <= 0 {
		conf.WriteTimeout = defaultWriteTimeout
	}
	if conf.DialTimeout <= 0 {
		conf.DialTimeout = defaultDialTimeout
	}
	if conf.MaxInFlight <= 0 {
		conf.MaxInFlight = defaultMaxInFlight
	}
	if conf.ThreadCount <= 0 {
		conf.ThreadCount = runtime.NumCPU()
	}
}
