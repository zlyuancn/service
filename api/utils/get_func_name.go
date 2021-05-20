package utils

import (
	"reflect"
	"runtime"
)

// 获取函数名
func GetFuncName(f interface{}) string {
	pc := reflect.ValueOf(f).Pointer()
	name := runtime.FuncForPC(pc).Name()
	return name
}
