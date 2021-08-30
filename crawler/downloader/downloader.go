package downloader

import (
	zapp_core "github.com/zly-app/zapp/core"

	"github.com/zly-app/service/crawler/core"
)

type Downloader struct {
}

func NewDownloader(app zapp_core.IApp) core.IDownloader {
	return &Downloader{}
}

func (d *Downloader) Close() error {
	return nil
}
