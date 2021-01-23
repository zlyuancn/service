
# mysql-binlog服务插件

> 提供用于 https://github.com/zly-app/zapp 的服务插件

# 说明

```text
mysql_binlog.RegistryService()              # 注册服务
mysql_binlog.WithMysqlBinlogService()       # 启用服务
mysql_binlog.RegistryMysqlBinlogHandler(...)    # 服务注入(注册记录事件handler)
```

# 示例

```go
// 注册服务
mysql_binlog.RegistryService()
// 启用服务
app := zapp.NewApp("test", mysql_binlog.WithMysqlBinlogService())
// 服务注入
mysql_binlog.RegistryMysqlBinlogHandler(app, &mysql_binlog.BaseEventHandler{})
// 运行
app.Run()
```

# 配置

> 默认服务类型为 `mysql-binlog`

```toml
[services.mysql-binlog]
# mysql 主机地址
Host="localhost:3306"
# 用户名, 最好是root
UserName="root"
# 密码
Password="yourpwd"
# 字符集, 一般为utf8mb4
Charset="utf8mb4"
# 包含的表正则匹配, 匹配的数据为 dbName.tableName
IncludeTableRegex=["^dbname\\.tablename$"]
# 排除的表正则匹配, 匹配的数据为 dbName.tableName
ExcludeTableRegex=[]
# 放弃没有表元数据的row事件
DiscardNoMetaRowEvent=true
# mysqldump执行路径, 如果为空则忽略mysqldump只使用binlog, mysqldump执行路径一般为mysqldump
DumpExecutionPath=""
# 忽略wkb数据解析错误, 一般为POINT, GEOMETRY类型
IgnoreWKBDataParseError=true
```
