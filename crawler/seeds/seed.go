package seeds

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/zly-app/service/crawler/core"
)

// 创建seed
func newSeed() *core.Seed {
	return &core.Seed{}
}

// 从原始数据生成seed
func MakeSeedOfRaw(raw string) (*core.Seed, error) {
	seed := newSeed()
	err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(raw, seed)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

/*
**创建种子
 uri 抓取连接
 parserMethod 解析方法, 可以是方法名或方法实体
*/
func NewSeed(uri string, parserMethod interface{}) *core.Seed {
	return newSeed().WithUri(uri).WithParserMethod(parserMethod)
}

// 将seed编码
func EncodeSeed(seed *core.Seed) (string, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(seed)
}
