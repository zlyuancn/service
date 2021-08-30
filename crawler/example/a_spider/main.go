package main

import (
	"github.com/zly-app/zapp"

	"github.com/zly-app/service/crawler"
)

type Spider struct {
	crawler.Spider
}

func main() {
	app := zapp.NewApp("a_spider", crawler.WithService())
	crawler.RegistrySpider(new(Spider))
	app.Run()
}
