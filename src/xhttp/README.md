## Mix XHTTP

A highly efficient HTTP library.

## Installation

```
go get github.com/mix-go/xhttp
```

## Functions

| Function                                                                                            | Description                       |  
|-----------------------------------------------------------------------------------------------------|-----------------------------------|
| xhttp.Fetch(ctx context.Context, method string, u string, opts ...RequestOption) (*Response, error) | Send a http request.              |
| xhttp.NewRequest(method string, u string, opts ...RequestOption) (*Request, error)                  | Create a request object.          |
| xhttp.Do(ctx context.Context, req *Request, opts ...RequestOption) (*Response, error)               | Send a http request.              |
| xhttp.DoRequest(ctx context.Context, req *http.Request, opts ...RequestOption) (*Response, error)   | Send a standard http request.     |
| xhttp.WithBody(body Body) RequestOption                                                             | Configure an option.              |
| xhttp.WithHeader(header http.Header) RequestOption                                                  | Configure an option.              |
| xhttp.WithContentType(contentType string) RequestOption                                             | Configure an option.              |
| xhttp.WithTimeout(timeout time.Duration) RequestOption                                              | Configure an option.              |
| xhttp.WithDebugFunc(f DebugFunc) RequestOption                                                      | Configure an option.              |
| xhttp.WithRetry(f RetryIfFunc, opts ...retry.Option) RequestOption                                  | Configure an option.              |
| xhttp.WithMiddleware(middlewares ...Middleware) RequestOption                                       | Configure an option.              |
| xhttp.BuildJSON(v interface{}) Body                                                                 | Generate json string.             |
| xhttp.BuildQuery(m map[string]string) Body                                                          | Generate urlencoded query string. |
| xhttp.Shutdown(ctx context.Context)                                                                 | Do shutdown.                      |

## Send a request

Send a simple request.

```go
url := "https://github.com/"
resp, err := xhttp.Fetch(context.Background(), http.MethodGet, url)
```

Configure some options.

```go
url := "https://github.com/"
header := make(http.Header)
header.Set("Authorization", "Bearer ***")
resp, err := xhttp.Fetch(context.Background(), http.MethodGet, url,
    xhttp.WithContentType("application/json"),
    xhttp.WithHeader(header))
```

Create a request object and send the request.

```go
url := "https://github.com/"
req, err := xhttp.NewRequest(http.MethodGet, url)
if err != nil {
    log.Fatal(err)
}
resp, err := xhttp.Do(context.Background(), req)
```

Create a standard request object and send the request.

```go
url := "https://github.com/"
req, err := http.NewRequest(http.MethodGet, url, nil)
if err != nil {
    log.Fatal(err)
}
resp, err := xhttp.DoRequest(context.Background(), req)
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
xhttp.Fetch(context.Background(), http.MethodPost, url, xhttp.WithDebugFunc(f))
```

The log object contains the following fields

```go
type Log struct {
    Context  context.Context `json:"context"`
    Duration time.Duration   `json:"duration"`
    Request  *Request        `json:"request"`  // The Request.RetryAttempts field records the number of retry attempts
    Response *Response       `json:"response"` // If request error this field is equal to nil
    Error    error           `json:"error"`
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
resp, err := xhttp.Fetch(context.Background(), http.MethodGet, url, xhttp.WithRetry(retryIf, retry.Attempts(2)))
```

Network error, no retry.

```go
url := "https://aaaaa.com/"
retryIf := func(resp *xhttp.Response, err error) error {
    if err != nil {
        return errors.Join(err, xhttp.ErrAbortRetry)
    }
    if resp.StatusCode != 200 {
        return fmt.Errorf("invalid status_code: %d", resp.StatusCode)
    }
    return nil
}
resp, err := xhttp.Fetch(context.Background(), http.MethodGet, url, xhttp.WithRetry(retryIf, retry.Attempts(2)))
```

## Middleware

Middleware configuration before or after.

```go
logicMiddleware := func(next xhttp.HandlerFunc) xhttp.HandlerFunc {
    return func(req *xhttp.Request) (*xhttp.Response, error) {
        // Before-logic
        fmt.Printf("Before: %s %s\n", req.Method, req.URL)

        // Call the next handler
        resp, err := next(req)

        // After-logic
        fmt.Printf("After: %s %s\n", req.Method, req.URL)

        return resp, err
    }
}
resp, err := xhttp.Fetch(context.Background(), http.MethodGet, "https://github.com/", xhttp.WithMiddleware(logicMiddleware))
```

## Shutdown

Before shutdown, all requests will be completed and work with middleware to save the response results to the database.

```go
ch := make(chan os.Signal)
signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
go func() {
    <-ch
    xhttp.Shutdown(context.Background())
}()
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
