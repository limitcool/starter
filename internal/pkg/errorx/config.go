package errorx

// 默认最大堆栈帧数
var defaultMaxStackFrames = 10

// SetMaxStackFrames 设置最大堆栈帧数
func SetMaxStackFrames(maxFrames int) {
	if maxFrames > 0 {
		defaultMaxStackFrames = maxFrames
	}
}

// GetMaxStackFrames 获取最大堆栈帧数
func GetMaxStackFrames() int {
	return defaultMaxStackFrames
}
