package crawler

import (
	"github.com/zly-app/service/crawler/core"
)

var _ core.ISpider = (*Spider)(nil)

type Spider struct {
}

func (s *Spider) Init() error {
	return nil
}

func (s *Spider) SubmitInitialSeed() []core.ISeed {
	return nil
}

func (s *Spider) Stop() error {
	return nil
}
