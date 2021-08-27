package core

// 队列
type IQueue interface {
	/*
		**将种子放入队列
		seed 种子
		back 是否放在后面
	*/
	Put(seed ISeed, back bool) error
	/*
	** 弹出一个种子
	 */
	Pop() (ISeed, error)
}
