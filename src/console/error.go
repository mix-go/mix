package console

import (
    "errors"
    "github.com/astaxie/beego/logs"
    "github.com/sirupsen/logrus"
    "runtime/debug"
)

const (
    // logrus
    LogrusType int = iota
    // beego
    BeegoType
)

// Error
type Error struct {
    LoggerManager loggerManager
    Dispatcher    string
}

// Handle
func (t *Error) Handle(err interface{}) {
    t.LoggerManager.Logf("%s", err)
}

type loggerManager struct {
    UseType      int
    LogrusLogger *logrus.Logger
    BeegoLogger  *logs.BeeLogger
}

func (t *loggerManager) AddLogrus(logger *logrus.Logger) {
    t.UseType = LogrusType
    t.LogrusLogger = logger
}

func (t *loggerManager) AddBeego(logger *logs.BeeLogger) {
    t.UseType = BeegoType
    t.BeegoLogger = logger
}

// 打印日志
func (t *loggerManager) Logf(format string, args ...interface{}) {
    format = format + "\n%s"
    args = append(args, string(debug.Stack()))
    switch t.UseType {
    case LogrusType:
        t.LogrusLogger.Errorf(format, args...)
        break
    case BeegoType:
        t.BeegoLogger.Error(format, args...)
        break
    }
}

// 创建 Error
func NewError(logger interface{}) *Error {
    loggerManager := loggerManager{}

    switch logger.(type) {
    case *logrus.Logger:
        loggerManager.AddLogrus(logger.(*logrus.Logger))
        break
    case *logs.BeeLogger:
        loggerManager.AddBeego(logger.(*logs.BeeLogger))
        break
    default:
        panic(errors.New("Unsupported logger type"))
    }

    return &Error{
        LoggerManager: loggerManager,
    }
}
