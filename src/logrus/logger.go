package logrus

import (
    "fmt"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    l "github.com/sirupsen/logrus"
    "io"
)

type Logger struct {
    *l.Logger
}

func (t *Logger) ErrorStack(err interface{}, stack string) {
    t.Logger.Errorf(fmt.Sprintf("%s\n%s", err, stack))
}

func NewLogger() *Logger {
    logger := l.New()

    formatter := new(l.TextFormatter)
    formatter.TimestampFormat = "2006-01-02 15:04:05"
    formatter.DisableQuote = true // 不转义换行符，为了保存错误堆栈到日志文件
    logger.Formatter = formatter

    return &Logger{logger}
}

func NewFileWriter(filename string, maxFiles int) io.Writer {
    writer, err := rotatelogs.New(
        filename+".%Y%m%d",
        rotatelogs.WithMaxAge(-1),
        rotatelogs.WithRotationCount(uint(maxFiles)),
    )
    if err != nil {
        panic(err)
    }
    return writer
}
