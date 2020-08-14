package console

import "github.com/mix-go/event"

type NotFoundError error
type UnsupportError error

// Error
type Error struct {
    Logger     Logger
    Dispatcher *event.EventDispatcher
}

type Logger interface {
    ErrorStack(err interface{}, stack string)
}

type HandleErrorEvent struct {
    event.EventTrait
    Error interface{}
}

// Handle
func (t *Error) Handle(err interface{}, stack []byte) {
    // dispatch
    t.dispatch(err)

    // log
    t.Logger.ErrorStack(err, string(stack))
}

func (t *Error) dispatch(err interface{}) {
    if t.Dispatcher == nil {
        return
    }
    eve := &HandleErrorEvent{
        Error: err,
    }
    t.Dispatcher.Dispatch(eve)
}

// 创建 Error
func NewError(logger Logger) *Error {
    return &Error{
        Logger: logger,
    }
}
