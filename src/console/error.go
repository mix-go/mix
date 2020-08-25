package console

import (
    event2 "github.com/mix-go/console/event"
    "github.com/mix-go/event"
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

type ErrorHandler struct {
    Logger     Logger
    Dispatcher *event.Dispatcher
}

func (t *ErrorHandler) Handle(err interface{}, stack []byte) {
    // dispatch
    t.dispatch(err)

    // log
    t.Logger.ErrorStack(err, string(stack))
}

func (t *ErrorHandler) dispatch(err interface{}) {
    if t.Dispatcher == nil {
        return
    }
    e := &event2.HandleErrorEvent{
        Error: err,
    }
    t.Dispatcher.Dispatch(e)
}

func NewError(logger Logger) Error {
    return &ErrorHandler{
        Logger: logger,
    }
}

type Error interface {
    Handle(err interface{}, stack []byte)
}

type Logger interface {
    ErrorStack(err interface{}, stack string)
}
