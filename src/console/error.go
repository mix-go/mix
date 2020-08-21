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

type Error struct {
    Logger     Logger
    Dispatcher *event.Dispatcher
}

type Logger interface {
    ErrorStack(err interface{}, stack string)
}

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
    e := &event2.HandleErrorEvent{
        Error: err,
    }
    t.Dispatcher.Dispatch(e)
}

func NewError(logger Logger) *Error {
    return &Error{
        Logger: logger,
    }
}
