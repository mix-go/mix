> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix Event

参考 PHP PSR-14 标准的事件调度库

Event dispatch library refer to PHP PSR-14 standard

> 该库还有 php 版本：https://github.com/mix-php/event

## Overview

事件调度是一种常见机制，允许开发人员轻松，一致地将逻辑注入应用程序，这在 PHP 中非常常见，于是我打造了这个 Go 版本的事件调度库，整体实现基于 [PSR-14](https://www.php-fig.org/psr/psr-14/)：

- 事件(Event)：一个自定义结构体，事件数据的载体
- 监听器(Listener)：负责处理对应的事件数据
- 事件调度器(Dispatcher)：用于触发某个事件

## Installation

- 安装

```
go get -u github.com/mix-go/event
```

## Usage

定义事件 `Event`，事件可以为任意结构体，只需要继承 `event.EventTrait` 即可

```
type CommandBeforeExecuteEvent struct {
    event.EventTrait
    Command interface{}
}
```

定义监听器 `Listener`，监听器是用户编写处理事件逻辑代码的地方，`Events` 方法返回一个要监听的事件类的数组，当这些事件触发时，会调用 `Process` 方法：

```
type CommandListener struct {
}

func (t *CommandListener) Events() []event.Event {
    return []event.Event{
        &event2.CommandBeforeExecuteEvent{},
    }
}

func (t *CommandListener) Process(e event.Event) {
    switch e.(type) {
    case *event2.CommandBeforeExecuteEvent:
        // 初始化全局对象
        globals.Init()

        // 设置守护
        if flag.Match("d", "daemon").Bool() {
            process.Daemon()
        }
        break
    }
}
```

触发某个事件 `Dispatcher`

```
e := &CommandBeforeExecuteEvent{
    Command: "foo",
}
d := event.NewDispatcher(&CommandListener)
d.Dispatch(e)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
