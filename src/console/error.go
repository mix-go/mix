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
    Handle(err interface{})
}

type Logger interface {
    ErrorStack(err interface{}, stack *[]byte)
}

type ErrorHandler struct {
    Logger     Logger
    Dispatcher event.Dispatcher
}

func (t *ErrorHandler) Handle(err interface{}) {
    LastError = err

    // dispatch
    t.dispatch(err)

    // log
    trace := debug.Stack()
    t.Logger.ErrorStack(err, &trace)
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
