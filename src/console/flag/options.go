package flag

var opts options

type options map[string]string

// Map 返回map
func (t *options) Map() map[string]string {
	return *t
}

// Options 获取全部命令行选项
func Options() *options {
	return &opts
}
