package xcli

// Command
type Command struct {
	// 命令名称
	Name string

	// 简短描述
	Short string

	// 详细描述
	Long string

	// 使用范例
	// 子命令："Usage: %s %s [ARG...]"
	// 单命令："Usage: %s [ARG...]"
	UsageF string

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
