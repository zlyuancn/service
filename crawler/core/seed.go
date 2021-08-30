package core

// 种子
type ISeed interface {
	// 原始数据
	Raw() string
	// 编码用于储存到队列
	Encode() (string, error)
}
