package logrus

import (
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "io"
)

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
