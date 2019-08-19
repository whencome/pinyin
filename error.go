package pinyin

import "errors"

var (
	// 尚未初始化
	ErrNotInitialized = errors.New("not yet initialized")
	// 已经读取到结尾
	ErrEOF = errors.New("end of file or string")
)
