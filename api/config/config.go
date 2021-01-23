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
	ShowApiResultInDevelop        bool   // 在开发环境中显示api结果
	ShowDetailedErrorInProduction bool   // 在生产环境显示详细的错误
}
