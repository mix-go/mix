package logrus

import (
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func NewFileWriter(filename string, count uint, size int64) *rotatelogs.RotateLogs {
    writer, err := rotatelogs.New(
        filename+".%Y%m%d",
        rotatelogs.WithLinkName(filename),
        rotatelogs.WithMaxAge(-1),
        rotatelogs.WithRotationCount(count),
        rotatelogs.WithRotationSize(size),
    )
    if err != nil {
        panic(err)
    }
    return writer
}
