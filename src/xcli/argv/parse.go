package argv

// 初始化
func init() {
	Parse()
}

// Parse 解析命令行参数
func Parse(singleton ...bool) {
	var s bool
	switch len(singleton) {
	case 0:
		s = false
	default:
		s = singleton[0]
	}

	prog = newProgram()
	cmd = newCommand(s)
}
