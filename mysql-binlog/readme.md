
# mysql-binlog服务

> 提供用于 https://github.com/zly-app/zapp 的服务

# 说明

```text
mysql_binlog.WithService()          # 启用服务
mysql_binlog.RegistryHandler(...)   # 服务注入(注册记录事件handler)
```

# 示例

```go
package main

import (
	mysql_binlog "github.com/zly-app/service/mysql-binlog"
	"github.com/zly-app/zapp"
)

func main() {
    // 启用服务
    app := zapp.NewApp("test", mysql_binlog.WithService())
    // 服务注入
    mysql_binlog.RegistryHandler(&mysql_binlog.BaseEventHandler{})
    // 运行
    app.Run()
}
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

# `Record`说明

## Record 结构

```go
type Record struct {
	Action    string                 `json:"action"`
	Old       map[string]interface{} `json:"old"`
	New       map[string]interface{} `json:"new"`
	DbName    string                 `json:"db_name"`
	TableName string                 `json:"table_name"`
	Timestamp uint32                 `json:"timestamp"`
}
```

1. `Action`表示动作, 只会是`update`,`insert`,`delete`这些值中的一个.

   ```
   当为 update 时,  Old 字段值为原始数据, New 的值为新数据.
   当为 insert 时,  Old 的值为 nil, New 字段值为新增的数据.
   当为 delete 时,  Old 字段值为原始数据, New 的值为 nil.
   ```

2. `DbName`表示数据库名.
3. `TableName`表示表名.
4. `Timestamp`表示数据发送改变时的时间戳, 单位为秒
