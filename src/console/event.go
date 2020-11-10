package console

import (
	"github.com/mix-go/event"
)

// HandleErrorEvent 处理错误事件
type HandleErrorEvent struct {
	event.EventTrait
	Error interface{}
}

// CommandBeforeExecuteEvent 命令行前置事件
type CommandBeforeExecuteEvent struct {
	event.EventTrait
	Command Command
}
