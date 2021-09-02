package main

import (
	"fmt"

	"github.com/zly-app/zapp"

	"github.com/zly-app/service/crawler"
	"github.com/zly-app/service/crawler/core"
)

type Spider struct {
	crawler core.ICrawler
	crawler.Spider
}

func (s *Spider) Init(crawler core.ICrawler) error {
	s.crawler = crawler
	return nil
}

func (s *Spider) SubmitInitialSeed() error {
	seed := s.crawler.NewSeed("https://www.baidu.com/", s.Parser)
	if err := s.crawler.PutSeed(seed, true); err != nil {
		return err
	}
	return nil
}

func (s *Spider) Parser(seed *core.Seed) error {
	fmt.Println(seed.Raw)
	fmt.Println(string(seed.HttpResponseBody))
	return nil
}

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService())
	crawler.RegistrySpider(new(Spider))
	app.Run()
}
