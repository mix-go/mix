package console

import (
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
    Dispatcher event.Dispatcher
}

func (t *ErrorHandler) Handle(err interface{}, stack []byte) {
    // dispatch
    t.dispatch(err)

    // log
    t.Logger.ErrorStack(err, stack)
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

type Error interface {
    Handle(err interface{}, stack []byte)
}

type Logger interface {
    ErrorStack(err interface{}, stack []byte)
}
