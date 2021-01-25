/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/25
   Description :
-------------------------------------------------
*/

package nsq_consume

const DefaultHeartbeatInterval = 30000 // 默认心跳间隔
const DefaultTimeout = 30000           // 默认超时

type Config struct {
	NsqdAddress       string // nsqd地址, localhost1:4150,localhost2:4150
	NsqLookupdAddress string // nsq发现服务地址, 优先级高于NsqdAddress, localhost1:4161,localhost2:4161
	AuthSecret        string // 验证秘钥
	HeartbeatInterval int64  // 心跳间隔(毫秒), 不能超过ReadTimeout, 0表示无
	ReadTimeout       int64  // 超时(毫秒)
	WriteTimeout      int64  // 超时(毫秒)
	DialTimeout       int64  // 超时(毫秒)
}

func newConfig() *Config {
	return &Config{
		NsqdAddress:       "",
		NsqLookupdAddress: "",
		AuthSecret:        "",
		HeartbeatInterval: DefaultHeartbeatInterval,
		ReadTimeout:       DefaultTimeout,
		WriteTimeout:      DefaultTimeout,
		DialTimeout:       DefaultTimeout,
	}
}

func (conf *Config) Check() {
	if conf.HeartbeatInterval <= 0 {
		conf.HeartbeatInterval = -1
	}
}
