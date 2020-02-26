package console

import (
    "fmt"
    "mix/src/bean"
    "mix/src/console/cli"
    "mix/src/console/flag"
    "reflect"
)

// 全局App
var App *Application

// 上下文
func Context() *bean.ApplicationContext {
    return App.Context
}

// 创建App
func NewApplication(definition ApplicationDefinition) *Application {
    app := &Application{
        ApplicationDefinition: definition,
    }
    app.Init()

    // 保存指针
    App = app

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
    // 依赖配置
    Beans []bean.BeanDefinition
    // 单命令
    // 该字段赋值后将忽略 Commands 字段的处理
    Command CommandDefinition
    // 命令集合
    Commands []CommandDefinition
}

// App
type Application struct {
    // App 定义
    ApplicationDefinition
    // 基础路径
    BasePath string
    // 应用上下文
    Context *bean.ApplicationContext
}

// 命令定义
type CommandDefinition struct {
    Name    string
    Usage   string
    Options []CommandOption
    Reflect func() reflect.Value
}

// 命令选项
type CommandOption struct {
    Names []string
    Usage string
}

// 初始化
func (p *Application) Init() {
    p.Context = bean.NewApplicationContext(p.Beans)
    p.BasePath = cli.Program.Dir
}

// 执行
func (p *Application) Run() {
    fmt.Println(cli.Command)
    fmt.Println(cli.Program)
    fmt.Println(flag.Options)
}
