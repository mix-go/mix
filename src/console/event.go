package console

import (
	"github.com/mix-go/event"
)

// 处理错误事件
type HandleErrorEvent struct {
	event.EventTrait
	Error interface{}
}

// 命令行前置事件
type CommandBeforeExecuteEvent struct {
	event.EventTrait
	Command Command
}
