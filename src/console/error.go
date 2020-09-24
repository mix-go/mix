package console

import (
    "github.com/mix-go/event"
    "runtime/debug"
)

type NotFoundError struct {
    error
}

func NewNotFoundError(err error) *NotFoundError {
    return &NotFoundError{err}
}

type UnsupportError struct {
    error
}

func NewUnsupportError(err error) *UnsupportError {
    return &UnsupportError{err}
}

type Error interface {
    Handle(err interface{}, trace ...bool)
}

type Logger interface {
    ErrorStack(err interface{}, stack *[]byte)
}

type ErrorHandler struct {
    Logger     Logger
    Dispatcher event.Dispatcher
}

func (t *ErrorHandler) Handle(err interface{}, trace ...bool) {
    LastError = err

    // dispatch
    t.dispatch(err)

    // log
    first := false
    if len(trace) > 0 {
        first = trace[0]
    }
    if first {
        trace := debug.Stack()
        t.Logger.ErrorStack(err, &trace)
    } else {
        t.Logger.ErrorStack(err, nil)
    }
}

func (t *ErrorHandler) dispatch(err interface{}) {
    if t.Dispatcher == nil {
        return
    }
    e := &HandleErrorEvent{
        Error: err,
    }
    t.Dispatcher.Dispatch(e)
}

func NewError(logger Logger) Error {
    return &ErrorHandler{
        Logger: logger,
    }
}
