package core

// 下载器
type IDownloader interface {
	// 关闭
	Close() error
}

// 代理
type IProxy interface {
}
