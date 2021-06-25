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
	RunF func()
	RunI RunI
	// 是否单命令
	Singleton bool
	// 是否为默认命令
	Default bool

	// handlers
	handlers []HandlerFunc
}

// AddOption
func (t *Command) AddOption(options ...*Option) *Command {
	t.Options = append(t.Options, options...)
	return t
}

// Use
func (t *Command) Use(h ...HandlerFunc) *Command {
	t.handlers = append(t.handlers, h...)
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
