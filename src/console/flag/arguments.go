package flag

var args arguments

type arguments []string

// Array 返回数组
func (t *arguments) Array() []string {
	return *t
}

// Values 返回值
func (t *arguments) Values() []*flagValue {
	args := *t
	var values []*flagValue
	for _, v := range args {
		values = append(values, &flagValue{v, true})
	}
	return values
}

// First 获取第一个参数
func (t *arguments) First() *flagValue {
	args := *t
	if len(args) == 0 {
		return &flagValue{}
	}
	return &flagValue{args[0], true}
}

// Arguments 获取全部命令行参数
func Arguments() *arguments {
	return &args
}
