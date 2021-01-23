/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package grpc

// 默认心跳时间
const defaultHeartbeatTime = 20000

// grpc服务配置
type Config struct {
	Bind          string // bind地址
	HeartbeatTime int    // 心跳时间(毫秒), 默认20000
}
