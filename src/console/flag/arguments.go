package flag

var args arguments

type arguments []string

// Array 返回数组
func (t *arguments) Array() []string {
	return *t
}

// Arguments 获取全部命令行参数
func Arguments() *arguments {
	return &args
}
