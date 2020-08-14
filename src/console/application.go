package console

import (
    "errors"
    "fmt"
    "github.com/mix-go/bean"
    "github.com/mix-go/console/argv"
    event2 "github.com/mix-go/console/event"
    "github.com/mix-go/console/flag"
    "github.com/mix-go/event"
    "reflect"
    "runtime/debug"
    "strings"
)

var (
    // 全局APP
    app *Application
    // 版本号
    Version = "1.0.0-beta6";
    // 最后的错误
    LastError interface{}
)

// App
func App() *Application {
    return app
}

// 上下文
func Context() *bean.ApplicationContext {
    return App().Context
}

// 创建App
func NewApplication(definition ApplicationDefinition, dispatcherName, errorName string) *Application {
    app = &Application{
        ApplicationDefinition: definition,
        DispatcherName:        dispatcherName,
        ErrorName:             errorName,
    }
    app.Init();
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
    // 命令集合
    Commands []CommandDefinition
}

// App
type Application struct {
    // App 定义
    ApplicationDefinition
    // Event Dispatcher
    DispatcherName string
    Dispatcher     *event.Dispatcher
    // Error
    ErrorName string
    Error     *Error
    // 基础路径
    BasePath string
    // 应用上下文
    Context *bean.ApplicationContext
    // 是否单命令
    Singleton bool
}

// 命令定义
type CommandDefinition struct {
    // 命令名称
    Name string
    // 使用描述
    Usage string
    // 选项
    Options []OptionDefinition
    // 反射
    Reflect func() reflect.Value
    // 是否单命令
    Singleton bool
}

// 命令选项
type OptionDefinition struct {
    Names []string
    Usage string
}

// 初始化
func (t *Application) Init() {
    t.Context = bean.NewApplicationContext(t.Beans)

    // 断言无法使用接口，由于没有泛型，导致这里 Dispatcher Error 无法实现 IoC
    // 等 go 推出泛型时再修改为接口
    t.Dispatcher = t.Context.Get(t.DispatcherName).(*event.Dispatcher)
    t.Error = t.Context.Get(t.ErrorName).(*Error)

    t.BasePath = argv.Program.Dir

    for _, c := range t.Commands {
        if c.Singleton {
            t.Singleton = true
            break
        }
    }
}

// 执行
func (t *Application) Run() {
    defer func() {
        if err := recover(); err != nil {
            LastError = err

            switch err.(type) {
            case NotFoundError, UnsupportError:
                fmt.Println(err)
                return
            }

            t.Error.Handle(err, debug.Stack())
        }
    }()

    if len(t.Commands) == 0 {
        panic(errors.New("Command cannot be empty"))
    }

    command := argv.Command
    if command == "" {
        if flag.BoolMatch([]string{"h", "help"}, false) {
            t.globalHelp()
            return
        }
        if flag.BoolMatch([]string{"v", "version"}, false) {
            t.version()
            return
        }

        options := flag.Options
        if len(options) == 0 {
            t.globalHelp()
            return
        } else if t.Singleton {
            t.call()
            return
        }

        f := ""
        for k, _ := range options {
            f = k
            break
        }
        p := argv.Program.Path
        panic(NotFoundError(errors.New(fmt.Sprintf("flag provided but not defined: '%s', see '%s --help'.", f, p))))
    } else if flag.Bool("help", false) {
        t.commandHelp()
        return
    }
    t.call()
}

// 执行命令
func (t *Application) call() {
    // 命令行选项效验
    t.validateOptions()

    // 提取命令
    var cmd *CommandDefinition
    command := argv.Command
    if t.Singleton {
        // 单命令
        for _, c := range t.Commands {
            if c.Singleton {
                cmd = &c
                break
            }
        }
        if cmd == nil {
            panic(errors.New("Singleton command not found"))
        }
    } else {
        for _, c := range t.Commands {
            if c.Name == command {
                cmd = &c
                break
            }
        }
    }
    if cmd == nil {
        panic(NotFoundError(errors.New(fmt.Sprintf("'%s' is not command, see '%s --help'.", command, argv.Program.Path))))
    }

    // 获取命令
    v := cmd.Reflect()

    // 触发执行命令前置事件
    e := &event2.CommandBeforeExecuteEvent{
        Command: v.Interface(),
    }
    t.Dispatcher.Dispatch(e)

    // 执行命令
    m := v.MethodByName("Main")
    if !m.IsValid() {
        panic(errors.New(fmt.Sprintf("'%s' Main method not found", fmt.Sprintf("%#v", v))))
    }
    m.Call([]reflect.Value{})
}

// 命令行选项效验
func (t *Application) validateOptions() {
    options := []OptionDefinition{}
    if !t.Singleton {
        for _, v := range t.Commands {
            if v.Name == argv.Command {
                options = v.Options
                break
            }
        }
    } else {
        for _, v := range t.Commands {
            if v.Singleton {
                options = v.Options
                break
            }
        }
    }
    if len(options) == 0 {
        return
    }

    flags := []string{}
    for _, o := range options {
        for _, v := range o.Names {
            if len(v) == 1 {
                flags = append(flags, fmt.Sprintf("-%s", v))
            } else {
                flags = append(flags, fmt.Sprintf("--%s", v))
            }
        }
    }
    inArray := func(value string, values []string) bool {
        for _, v := range values {
            if v == value {
                return true
            }
        }
        return false
    }
    for f, _ := range flag.Options {
        if !inArray(f, flags) {
            p := argv.Program.Path
            c := argv.Command
            if c != "" {
                c = fmt.Sprintf(" %s", c)
            }
            panic(NotFoundError(errors.New(fmt.Sprintf("flag provided but not defined: '%s', see '%s%s --help'.", f, p, c))))
        }
    }
}

// 全局帮助
func (t *Application) globalHelp() {
    program := argv.Program.Path
    fg := ""
    if !t.Singleton {
        fg = " [OPTIONS] COMMAND"
    }
    fmt.Println(fmt.Sprintf("Usage: %s%s [opt...]", program, fg))
    t.printGlobalOptions()
    if !t.Singleton {
        t.printCommands()
    } else {
        t.printCommandOptions()
    }
    fmt.Println("")
    fg = ""
    if !t.Singleton {
        fg = " COMMAND"
    }
    fmt.Println("")
    fmt.Println(fmt.Sprintf("Run '%s%s --help' for more information on a command.", program, fg))
    fmt.Println("")
    fmt.Println("Developed with Mix Go framework. (openmix.org/mix-go)")
}

// 命令帮助
func (t *Application) commandHelp() {
    program := argv.Program.Path
    command := argv.Command
    fmt.Println(fmt.Sprintf("Usage: %s %s [opt...]", program, command))
    t.printCommandOptions()
    fmt.Println("")
    fmt.Println("Developed with Mix Go framework. (openmix.org/mix-go)")
}

// 打印全局选项
func (t *Application) printGlobalOptions() {
    tabs := "\t";
    fmt.Println("");
    fmt.Println("Global Options:");
    fmt.Println(fmt.Sprintf("  -h, --help%sPrint usage", tabs));
    fmt.Println(fmt.Sprintf("  -v, --version%sPrint version information", tabs));
}

// 打印命令
func (t *Application) printCommands() {
    fmt.Println("");
    fmt.Println("Commands:");
    for _, v := range t.Commands {
        command := v.Name;
        usage := v.Usage
        fmt.Println(fmt.Sprintf("  %s\t%s", command, usage))
    }
}

// 打印命令选项
func (t *Application) printCommandOptions() {
    options := []OptionDefinition{}
    if !t.Singleton {
        for _, v := range t.Commands {
            if v.Name == argv.Command {
                options = v.Options
                break
            }
        }
    } else {
        for _, v := range t.Commands {
            if v.Singleton {
                options = v.Options
                break
            }
        }
    }
    if len(options) == 0 {
        return
    }

    fmt.Println("");
    fmt.Println("Command Options:");
    for _, o := range options {
        flags := []string{}
        for _, v := range o.Names {
            if len(v) == 1 {
                flags = append(flags, fmt.Sprintf("-%s", v))
            } else {
                flags = append(flags, fmt.Sprintf("--%s", v))
            }
        }
        fg := strings.Join(flags, ", ")
        usage := o.Usage
        fmt.Println(fmt.Sprintf("  %s\t%s", fg, usage));
    }
}

// 版本号
func (t *Application) version() {
    appName := t.AppName;
    appVersion := t.AppVersion;
    frameworkVersion := Version;
    fmt.Println(fmt.Sprintf("%s version %s, framework version %s", appName, appVersion, frameworkVersion));
}
