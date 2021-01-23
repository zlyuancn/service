/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/23
   Description :
-------------------------------------------------
*/

package mysql_binlog

// MysqlBinlogService配置
type Config struct {
	Host                    string   // mysql 主机地址
	UserName                string   // 用户名, 最好是root
	Password                string   // 密码
	Charset                 *string  // 字符集, 一般为utf8mb4
	IncludeTableRegex       []string // 包含的表正则匹配, 匹配的数据为 dbName.tableName
	ExcludeTableRegex       []string // 排除的表正则匹配, 匹配的数据为 dbName.tableName
	DiscardNoMetaRowEvent   bool     // 放弃没有表元数据的row事件
	DumpExecutionPath       string   // mysqldump执行路径, 如果为空则忽略mysqldump只使用binlog, mysqldump执行路径一般为mysqldump
	IgnoreWKBDataParseError bool     // 忽略wkb数据解析错误, 一般为POINT, GEOMETRY类型
}
