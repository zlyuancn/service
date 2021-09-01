package response_middleware

import (
	"fmt"

	"github.com/zly-app/service/crawler/config"
	"github.com/zly-app/service/crawler/core"
)

type CheckHttpStatusCode struct {
	core.MiddlewareBase
}

func NewCheckSeedIsValidMiddleware() core.IRequestMiddleware {
	return new(CheckHttpStatusCode)
}

func (m *CheckHttpStatusCode) Name() string { return "CheckHttpStatusCode" }
func (m *CheckHttpStatusCode) Process(crawler core.ICrawler, seed *core.Seed) (*core.Seed, error) {
	if seed.Response == nil {
		return seed, nil
	}

	// 检查期望值
	if len(config.Conf.Spider.ExpectHttpStatusCode) > 0 {
		for _, expect := range config.Conf.Spider.ExpectHttpStatusCode {
			if expect == seed.Response.StatusCode {
				return seed, nil
			}
		}
		return nil, fmt.Errorf("收到非期望的http状态码: %d", seed.Response.StatusCode)
	}

	// 检查排除值
	for _, exclude := range config.Conf.Spider.InvalidHttpStatusCode {
		if exclude == seed.Response.StatusCode {
			return nil, fmt.Errorf("收到无效的http状态码: %d", seed.Response.StatusCode)
		}
	}

	return seed, nil
}
