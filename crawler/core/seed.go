package core

import (
	"fmt"
	"net/http"

	"github.com/zly-app/service/crawler/utils"
)

// 种子数据
type Seed struct {
	// 原始数据, 构建种子时的数据
	Raw string `json:"-"`

	Request *http.Request `json:"-"`
	// 响应, 注意: 不能使用它的body, 而是应该使用 ResponseBody
	Response *http.Response `json:"-"`
	// 响应数据
	ResponseBody []byte `json:"-"`

	// uri
	Uri string
	// 解析方法名
	ParserMethod string
	// 检查期望方法名
	CheckExpectMethod string
}

// 设置Uri
func (s *Seed) WithUri(uri string) *Seed {
	s.Uri = uri
	return s
}

// 设置解析方法
func (s *Seed) WithParserMethod(parserMethod interface{}) *Seed {
	switch t := parserMethod.(type) {
	case string:
		s.ParserMethod = t
	case ParserMethod:
		s.ParserMethod = utils.Reflect.GetFuncName(t)
	default:
		panic(fmt.Errorf("无法获取方法名: [%T]%v", parserMethod, parserMethod))
	}
	return s
}

// 设置检查期望响应方法
func (s *Seed) WithCheckExpectMethod(checkMethod interface{}) *Seed {
	switch t := checkMethod.(type) {
	case string:
		s.CheckExpectMethod = t
	case ParserMethod:
		s.CheckExpectMethod = utils.Reflect.GetFuncName(t)
	default:
		panic(fmt.Errorf("无法获取方法名: [%T]%v", checkMethod, checkMethod))
	}
	return s
}
