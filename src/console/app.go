package console

import (
	"mix/src/bean"
	"reflect"
)

// 创建App
func NewApplication(definition *ApplicationDefinition) *Application {
	app := &Application{
		ApplicationDefinition: definition,
	}
	app.Init()
	return app
}

// App 定义
type ApplicationDefinition struct {
	// 应用名称
	AppName string
	// 应用版本
	AppVersion string
	// 应用调试
	AppDebug bool
	// 基础路径
	BasePath string
	// 单命令
	// 该字段赋值后将忽略 Commands 字段的处理
	Command CommandDefinition
	// 命令集合
	Commands []CommandDefinition
	// 依赖配置
	Beans []bean.BeanDefinition
}

// App
type Application struct {
	// App 定义
	*ApplicationDefinition
	// 应用上下文
	Context *bean.ApplicationContext
}

// 命令定义
type CommandDefinition struct {
	Reflect     func() reflect.Value
	Description string
	Options     []CommandOption
}

// 命令选项
type CommandOption struct {
	Names       []string
	Description string
}

// 初始化
func (p *Application) Init() {
	p.Context = bean.NewApplicationContext(p.Beans)
}
