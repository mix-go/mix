> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix XUtil

A set of tools that keep Golang sweet.

## Installation

```
go get github.com/mix-go/xutil
```

## xhttp

| Function                                                                                   | Description             |  
|--------------------------------------------------------------------------------------------|-------------------------|
| xhttp.Request(method string, u string, body string, opts ...Options) ([]byte, error)       | Execute a http request. |
| xhttp.RequestString(method string, u string, body string, opts ...Options) (string, error) | Execute a http request. |

## xslices

| Function                                              | Description                                          |  
|-------------------------------------------------------|------------------------------------------------------|
| xslices.InArray[T comparable](item T, slice []T) bool | Searches if the specified value exists in the array. |

## xstrings

| Function                          | Description                                                    |  
|-----------------------------------|----------------------------------------------------------------|
| xstrings.IsNumeric(s string) bool | Used to check if the variable is a number or a numeric string. |

## xfmt

A formatting library that ca print data inside nested pointer addresses of
structures, [see more](https://github.com/mix-go/mix/blob/master/src/xutil/xfmt/README.md).

The supported methods are identical to the `fmt` system library

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

Environment configuration library with type
conversion, [see more](https://github.com/mix-go/mix/blob/master/src/xutil/xenv/README.md).

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
