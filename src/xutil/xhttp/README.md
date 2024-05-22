> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix XHttp

A highly efficient HTTP library.

## Installation

```
go get github.com/mix-go/xutil
```

## Functions

| Function                                                                         | Description                      |  
|----------------------------------------------------------------------------------|----------------------------------|
| xhttp.Request(method string, u string, opts ...RequestOption) (*Response, error) | Execute an http request.         |
| xhttp.Do(req *http.Request, opts ...RequestOption) (*Response, error)            | Execute an http request.         |
| xhttp.WithBody(body Body) RequestOption                                          | Set configuration item           |
| xhttp.WithHeader(header http.Header) RequestOption                               | Set configuration item           |
| xhttp.WithContentType(contentType string) RequestOption                          | Set configuration item           |
| xhttp.WithTimeout(timeout time.Duration) RequestOption                           | Set configuration item           |
| xhttp.WithDebugFunc(f DebugFunc) RequestOption                                   | Set configuration item           |
| xhttp.WithRetry(f RetryIfFunc, opts ...retry.Option) RequestOption               | Set configuration item           |
| xhttp.BuildJSON(v interface{}) Body                                              | Generate json string             |
| xhttp.BuildQuery(m map[string]string) Body                                       | Generate urlencoded query string |

## Debug Log

By configuring `DebugFunc`, you can use any logging library to print log information here.

- Global configuration

```go
xhttp.DefaultOptions.DebugFunc = func(l *Log) {
    log.Println(l)
}
```

- Single request configuration

```go
f := func(l *Log) {
    log.Println(l)
}
xhttp.Request("POST", url, xhttp.WithDebugFunc(f))
```

The log object contains the following fields

```go
type Log struct {
	Duration time.Duration `json:"duration"`
	Request  *XRequest     `json:"request"`  // The XRequest.RetryAttempts field records the number of retry attempts
	Response *XResponse    `json:"response"` // If request error this field is equal to nil
	Error    error         `json:"error"`
}
```

## Retry

Set the conditions for determining retries, and specify various options such as the number of attempts.

```go
url := "https://aaaaa.com/"
retryIf := func(resp *xhttp.Response, err error) error {
    if err != nil {
        return err
    }
    if resp.StatusCode != 200 {
        return fmt.Errorf("invalid status_code: %d", resp.StatusCode)
    }
    return nil
}
resp, err := xhttp.Request("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(2)))
```

Network error, no retry.

```go
url := "https://aaaaa.com/"
retryIf := func(resp *xhttp.XResponse, err error) error {
    if err != nil {
        return errors.Join(err, xhttp.ErrAbortRetry)
    }
    if resp.StatusCode != 200 {
        return fmt.Errorf("invalid status_code: %d", resp.StatusCode)
    }
    return nil
}
resp, err := xhttp.Request("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(2)))
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
