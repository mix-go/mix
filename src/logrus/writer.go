package logrus

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// 创建文件日志处理
// count, size == 0 时不轮转
func NewFileWriter(filename string, count uint, size int64) *rotatelogs.RotateLogs {
	options := []rotatelogs.Option{
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(-1),
	}
	if count > 0 {
		options = append(options, rotatelogs.WithRotationCount(count))
	}
	if size > 0 {
		options = append(options, rotatelogs.WithRotationSize(size))
	}
	writer, err := rotatelogs.New(filename+".%Y%m%d", options...)
	if err != nil {
		panic(err)
	}
	return writer
}
