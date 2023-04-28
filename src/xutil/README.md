> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XUtil

一套让 Golang 保持甜美的工具。

A set of tools that keep Golang sweet.

## Installation

```
go get github.com/mix-go/xutil
```

## List of functions

执行一个http请求。

Execute an http request.

```go
http.Request(method string, u string, body string, opts ...Options) ([]byte, error)
```

```go
http.RequestString(method string, u string, body string, opts ...Options) (string, error)
```

搜索数组中是否存在指定的值。

Searches if the specified value exists in the array.

```go
slices.InArray[T comparable](item T, slice []T) bool
```

用于检测变量是否为数字或数字字符串。

Used to check if the variable is a number or a numeric string.

```go
strings.IsNumeric(s string) bool
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
