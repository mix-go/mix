package flag

var opts options

type options map[string]string

// Map {"--foo": "bar"}
func (t *options) Map() map[string]string {
	return *t
}

// Options 获取全部命令行选项
func Options() *options {
	return &opts
}
