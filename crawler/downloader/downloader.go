package downloader

import (
	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/service/crawler/core"
)

type Downloader struct {
	app zapp_core.IApp
}

func (d *Downloader) Download(seed *core.Seed) (*core.Seed, error) {
	if seed.Uri == "" {
		return seed, nil
	}

	return seed, nil
}

func (d *Downloader) Close() error {
	return nil
}

func NewDownloader(app zapp_core.IApp) core.IDownloader {
	return &Downloader{
		app: app,
	}
}
