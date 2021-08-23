/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/3
   Description :
-------------------------------------------------
*/

package mysql_binlog

import (
	jsoniter "github.com/json-iterator/go"
)

type Record struct {
	Action    string                 `json:"action"`
	Old       map[string]interface{} `json:"old"`
	New       map[string]interface{} `json:"new"`
	DbName    string                 `json:"db_name"`
	TableName string                 `json:"table_name"`
	Timestamp uint32                 `json:"timestamp"`
}

// 获取old数据的json输出
func (r *Record) OldString() string {
	text, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(r.Old)
	return text
}

// 获取new数据的json输出
func (r *Record) NewString() string {
	text, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(r.New)
	return text
}

// 将record的所有数据转为json
func (r *Record) String() string {
	text, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(r)
	return text
}

// 将old数据解析到a
func (r *Record) UnmarshalOld(a interface{}) error {
	return ScanMysqlData(r.Old, a)
}

// 将new数据解析到a
func (r *Record) UnmarshalNew(a interface{}) error {
	return ScanMysqlData(r.New, a)
}
