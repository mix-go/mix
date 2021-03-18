package xcli

// Command
type Command struct {
	// 命令名称
	Name string
	// 使用描述
	Usage string
	// 选项
	Options []*Option
	// 执行
	Run  func()
	RunI RunI
	// 是否单命令
	Singleton bool
	// 是否为默认命令
	Default bool
}

// AddOption
func (t *Command) AddOption(options ...*Option) *Command {
	t.Options = append(t.Options, options...)
	return t
}

// RunI
type RunI interface {
	Main()
}

// Option
type Option struct {
	Names []string
	Usage string
}
