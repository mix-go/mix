## Mix Event

参考 PHP PSR-14 标准的事件调度库

Event dispatch library refer to PHP PSR-14 standard

> 该库还有 php 版本：https://github.com/mix-php/event

## Overview

事件调度是一种常见机制，允许开发人员轻松，一致地将逻辑注入应用程序，这在 PHP 中非常常见，于是我打造了这个 Go 版本的事件调度库。

## Installation

- 安装

```
go get -u github.com/mix-go/event
```

## Usage

定义监听器

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

定义事件

```
type CommandBeforeExecuteEvent struct {
    event.EventTrait
    Command interface{}
}
```

调度

```
e := &CommandBeforeExecuteEvent{
    Command: "foo",
}
d := NewDispatcher(&CommandListener)
d.Dispatch(e)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
