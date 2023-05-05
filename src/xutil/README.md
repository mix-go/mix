> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XUtil

A set of tools that keep Golang sweet.

一套让 Golang 保持甜美的工具。

## Installation

```
go get github.com/mix-go/xutil
```

## xhttp

Execute an http request.

执行一个http请求。

### xhttp.Request

```go
xhttp.Request(method string, u string, body string, opts ...Options) ([]byte, error)
```

### xhttp.RequestString

```go
xhttp.RequestString(method string, u string, body string, opts ...Options) (string, error)
```

## xslices

### xslices.InArray

Searches if the specified value exists in the array.

搜索数组中是否存在指定的值。

```go
xslices.InArray[T comparable](item T, slice []T) bool
```

## xstrings

### xstrings.IsNumeric

Used to check if the variable is a number or a numeric string.

用于检测变量是否为数字或数字字符串。

```go
xstrings.IsNumeric(s string) bool
```

## xfmt

可以打印指针字段内部的数据结构，支持的方法与 `fmt` 系统库完全一致，[查看更多用法](xfmt/README.md)

- `Sprintf(format string, args ...interface{}) string`
- `Sprint(args ...interface{}) string`
- `Sprintln(args ...interface{}) string`
- `Printf(format string, args ...interface{})`
- `Print(args ...interface{})`
- `Println(args ...interface{})`

## xenv

具有类型转换功能的环境配置库，[查看更多用法](xenv/README.md)

获取环境变量并设置默认值

~~~go
i := xenv.Getenv("key").String("default")
i := xenv.Getenv("key").Bool(false)
i := xenv.Getenv("key").Int64(123)
i := xenv.Getenv("key").Float64(123.4)
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
