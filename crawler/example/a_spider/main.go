package main

import (
	"fmt"

	"github.com/zly-app/zapp"

	"github.com/zly-app/service/crawler"
	"github.com/zly-app/service/crawler/core"
	"github.com/zly-app/service/crawler/seeds"
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
	if err := s.crawler.PutSeed(seeds.NewSeed("", s.Parser), true); err != nil {
		return err
	}
	return nil
}

func (s *Spider) Parser(seed *core.Seed) error {
	fmt.Println(seed)
	return nil
}

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService())
	crawler.RegistrySpider(new(Spider))
	app.Run()
}
