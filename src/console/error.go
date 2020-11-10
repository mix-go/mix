package console

import (
	"github.com/mix-go/event"
	"runtime/debug"
)

// 未找到错误
type NotFoundError struct {
	error
}

// 创建未找到错误
func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}

// 不支持错误
type UnsupportError struct {
	error
}

// 创建不支持错误
func NewUnsupportError(err error) *UnsupportError {
	return &UnsupportError{err}
}

// 错误接口
type Error interface {
	Handle(err interface{})
}

// 日志接口
type Logger interface {
	ErrorStack(err interface{}, stack *[]byte)
}

// 错误处理结构体
type ErrorHandler struct {
	Logger     Logger
	Dispatcher event.Dispatcher
}

// 处理错误
func (t *ErrorHandler) Handle(err interface{}) {
	LastError = err

	// dispatch
	t.dispatch(err)

	// log
	trace := debug.Stack()
	t.Logger.ErrorStack(err, &trace)
}

// 调度事件
func (t *ErrorHandler) dispatch(err interface{}) {
	if t.Dispatcher == nil {
		return
	}
	e := &HandleErrorEvent{
		Error: err,
	}
	t.Dispatcher.Dispatch(e)
}

// 创建错误对象
func NewError(logger Logger) Error {
	return &ErrorHandler{
		Logger: logger,
	}
}
