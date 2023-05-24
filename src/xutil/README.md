> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XUtil

A set of tools that keep Golang sweet.

一套让 Golang 保持甜美的工具。

## Installation

```
go get github.com/mix-go/xutil
```

## xhttp

| Function                                                                                   | Description              |  
|--------------------------------------------------------------------------------------------|--------------------------|
| xhttp.Request(method string, u string, body string, opts ...Options) ([]byte, error)       | Execute an http request. |
| xhttp.RequestString(method string, u string, body string, opts ...Options) (string, error) | Execute an http request. |

## xslices

| Function                                              | Description                                          |  
|-------------------------------------------------------|------------------------------------------------------|
| xslices.InArray[T comparable](item T, slice []T) bool | Searches if the specified value exists in the array. |

## xstrings

| Function                          | Description                                                    |  
|-----------------------------------|----------------------------------------------------------------|
| xstrings.IsNumeric(s string) bool | Used to check if the variable is a number or a numeric string. |

## xfmt

可以打印结构体嵌套指针地址内部数据的格式化库，[查看更多](xfmt/README.md)。

支持的方法与 `fmt` 系统库完全一致

| Function                                                | Description       |  
|---------------------------------------------------------|-------------------|
| xfmt.Sprintf(format string, args ...interface{}) string |                   |
| xfmt.Sprint(args ...interface{}) string                 |                   |
| xfmt.Sprintln(args ...interface{}) string               |                   |
| xfmt.Printf(format string, args ...interface{})         |                   |
| xfmt.Print(args ...interface{})                         |                   |
| xfmt.Println(args ...interface{})                       |                   |
| xfmt.Disable()                                          | Equivalent to fmt |
| xfmt.Enable()                                           |                   |

## xenv

具有类型转换功能的环境配置库，[查看更多](xenv/README.md)。

| Function                                  | Description |  
|-------------------------------------------|-------------|
| err := xenv.Load(".env")                  |             |
| err := xenv.Overload(".env")              |             |
| i := xenv.Getenv("key").String("default") |             |
| i := xenv.Getenv("key").Bool(false)       |             |
| i := xenv.Getenv("key").Int64(123)        |             |
| i := xenv.Getenv("key").Float64(123.4)    |             |

## License

Apache License Version 2.0, http://www.apache.org/licenses/
