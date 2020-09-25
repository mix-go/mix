## Mix Logrus

基于 Logrus 扩展的日志库，支持行号、文件切分、调用堆栈

Based on Logrus extended log library, support file line, file rotate, call stack

## Overview

与原版 `Logrus` 有哪些不同？

- 开启了文件行号的支持，方便调试程序，使用 `logger.ReportCaller` 关闭。
- 设置了国内使用习惯的时间格式 `2006-01-02 15:04:05.000`，使用 `logger.Formatter.TimestampFormat` 修改。
- 扩展了日志文件轮转功能，并自动保留指定天数的文件数。
- 扩展了 `mix-go/console` 的 panic 捕获功能，可以实现捕获 panic 堆栈信息到日志，这个功能让我们不需要再从标准输出中查找 panic 错误信息。
- 支持 `GORM` 的 SQL 日志打印美化。

## Installation

- 安装

```
go get -u github.com/mix-go/logrus
```

## Usage

- 配置 `os.Stdout` 输出

~~~
logger := logrus.NewLogger()
logger.SetOutput(os.Stdout)
~~~

输出格式

~~~
INFO[2020-09-18 18:36:23.342]hello.go:16 This is the content
~~~

- 同时配置  `os.Stdout` 与 `io.Writer` 输出

>[info] NewFileWriter(filename string, maxFiles int) io.Writer

~~~
logger := logrus.NewLogger()
file := logrus.NewFileWriter("/tmp/logs/test.log", 7)
writer := io.MultiWriter(os.Stdout, file)
logger.SetOutput(writer)
~~~

输出格式

~~~
time=2020-09-18 16:18:51.470 level=info msg=This is the content file=hello.go:16
~~~

文件轮转格式

~~~
/tmp/logs/test.log.20200916
~~~

## GORM 支持

GORM 支持 `SetLogger` 打印日志，但是由于 GORM 是使用的 `Print` 方法，因此打印的多个参数会直接拼接到一起，难以阅读

```
time=2020-09-25 11:06:42.103 level=info msg=sql/foo.go:335319.978µsSELECT * FROM `foo`  WHERE _id = ? `foo`.`id` ASC LIMIT 1[122569]1 file=logger.go:32
```

开启 `GORM` 支持

```
logger := logrus.NewLogger()
logger.SupportGORM = true // 这个配置会美化日志

gorm.LogMode(true)
gorm.SetLogger(logger)
```

开启后日志

```
time=2020-09-25 11:06:42.103 level=info msg=sql /foo.go:335 319.978µs SELECT * FROM `foo`  WHERE _id = ? `foo`.`id` ASC LIMIT 1 [122569] 1 file=logger.go:32
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
