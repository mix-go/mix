package flag

var args arguments

type arguments []string

// Array 返回数组
func (t *arguments) Array() []string {
	return *t
}

// First 获取第一个参数
func (t *arguments) First() *flagValue {
	a := *t
	if len(a) == 0 {
		return &flagValue{}
	}
	return &flagValue{a[0], true}
}

// Arguments 获取全部命令行参数
func Arguments() *arguments {
	return &args
}
