package console

import (
    "github.com/mix-go/event"
)

type HandleErrorEvent struct {
    event.EventTrait
    Error interface{}
}

type CommandBeforeExecuteEvent struct {
    event.EventTrait
    Command Command
}
