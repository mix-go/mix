package console

import (
    "errors"
    "github.com/astaxie/beego/logs"
    "github.com/sirupsen/logrus"
)

const (
    // logrus
    LogrusType int = iota
    // beego
    BeegoType
)

// Error
type errorHandler struct {
    Level      int
    Logger     LoggerManager
    Dispatcher string
}

// Handle Exception
func (t *errorHandler) Handle(err interface{}) {
    t.Logger.logf("%s", err)
}

type Error interface {
    Handle(err interface{})
}

// LoggerManager
type LoggerManager struct {
    UseType      int
    LogrusLogger *logrus.Logger
    BeegoLogger  *logs.BeeLogger
}

func (t *LoggerManager) AddLogrus(logger *logrus.Logger) {
    t.UseType = LogrusType
    t.LogrusLogger = logger
}

func (t *LoggerManager) AddBeego(logger *logs.BeeLogger) {
    t.UseType = BeegoType
    t.BeegoLogger = logger
}

// 打印日志
func (t *LoggerManager) logf(format string, args ...interface{}) {
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
func NewError(level int, logger interface{}) Error {
    loggerManager := LoggerManager{}

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

    return &errorHandler{
        Level:  level,
        Logger: loggerManager,
    }
}
