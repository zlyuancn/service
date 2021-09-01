package crawler

import (
	"github.com/zly-app/service/crawler/core"
)

var _ core.ISpider = (*Spider)(nil)

type Spider struct {
}

func (s *Spider) Init(crawler core.ICrawler) error {
	return nil
}

func (s *Spider) SubmitInitialSeed() error {
	return nil
}

func (s *Spider) Stop() error {
	return nil
}
