/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/21
   Description :
-------------------------------------------------
*/

package config

var Conf Config

// api服务配置
type Config struct {
	Bind                          string // bind地址
	IPWithNginxForwarded          bool   // 适配nginx的Forwarded获取ip, 优先级高于nginx的Real
	IPWithNginxReal               bool   // 适配nginx的Real获取ip, 优先级高于sock连接的ip
	LogApiResultInDevelop         bool   // 在开发环境中输出api结果
	SendDetailedErrorInProduction bool   // 在生产环境发送详细的错误到客户端
}
