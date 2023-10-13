> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix XHttp

A highly efficient HTTP library.

## Installation

```
go get github.com/mix-go/xutil
```

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
	Request  *XRequest     `json:"request"`  // The XRequest.RetryCount field records the number of retries that have been completed.
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

## License

Apache License Version 2.0, http://www.apache.org/licenses/
