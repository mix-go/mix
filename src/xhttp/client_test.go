package xhttp_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mix-go/xhttp"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"sync"
	"testing"
)

func TestFetch(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.Fetch(context.Background(), http.MethodGet, url)

	a.Nil(err)
	a.Equal(resp.StatusCode, 200)
}

func TestDo(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	req, err := xhttp.NewRequest(http.MethodGet, url)
	a.Nil(err)
	resp, err := xhttp.Do(context.Background(), req)
	a.Nil(err)

	a.Equal(resp.StatusCode, 200)
}

func TestDoRequest(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	a.Nil(err)
	resp, err := xhttp.DoRequest(context.Background(), req)
	a.Nil(err)

	a.Equal(resp.StatusCode, 200)
}

func TestRequestPOST(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.Fetch(context.Background(), "POST", url, xhttp.WithBodyString("abc"), xhttp.WithContentType("application/json"))

	a.Nil(err)
	a.Equal(resp.StatusCode, 404)
}

func TestRequestError(t *testing.T) {
	a := assert.New(t)

	url := "https://aaaaa.com/"
	resp, err := xhttp.Fetch(context.Background(), "GET", url)

	a.NotNil(err)
	a.Nil(resp)
}

func TestDebugAndRetryFail(t *testing.T) {
	a := assert.New(t)

	count := 0
	xhttp.DefaultOptions.DebugFunc = func(l *xhttp.Log) {
		log.Printf("%+v %+v %+v %+v\n", l.Duration, l.Request, l.Response, l.Error)
		count++
	}

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
	resp, err := xhttp.Fetch(context.Background(), "GET", url, xhttp.WithRetry(retryIf, retry.Attempts(2)))

	a.NotNil(err)
	a.Nil(resp)
	a.Contains(err.Error(), "attempts fail")
	a.Equal(count, 2)
}

func TestDebugAndRetrySuccess(t *testing.T) {
	a := assert.New(t)

	count := 0
	xhttp.DefaultOptions.DebugFunc = func(l *xhttp.Log) {
		log.Printf("%+v %+v %+v %+v\n", l.Duration, l.Request, l.Response, l.Error)
		count++
	}

	url := "https://aaaaa.com/"
	retryIf := func(resp *xhttp.Response, err error) error {
		if count == 1 {
			return errors.New("the first request failed")
		}
		return nil
	}
	_, err := xhttp.Fetch(context.Background(), "GET", url, xhttp.WithRetry(retryIf, retry.Attempts(3)))

	a.Nil(err)
	a.Equal(count, 2)
}

func TestDebugAndRetryAbort(t *testing.T) {
	a := assert.New(t)

	count := 0
	xhttp.DefaultOptions.DebugFunc = func(l *xhttp.Log) {
		log.Printf("%+v %+v %+v %+v\n", l.Duration, l.Request, l.Response, l.Error)
		count++
	}

	url := "https://aaaaa.com/"
	retryIf := func(resp *xhttp.Response, err error) error {
		if err != nil {
			if count == 1 {
				return err
			}
			return errors.Join(err, xhttp.ErrAbortRetry)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status_code: %d", resp.StatusCode)
		}
		return nil
	}
	resp, err := xhttp.Fetch(context.Background(), "GET", url, xhttp.WithRetry(retryIf, retry.Attempts(3)))

	a.NotNil(err)
	a.Nil(resp)
	a.Contains(err.Error(), "xhttp: abort further retries")
	a.Equal(count, 2)
}

func TestMiddlewares(t *testing.T) {
	a := assert.New(t)

	logicMiddleware := func(next xhttp.HandlerFunc) xhttp.HandlerFunc {
		return func(req *xhttp.Request, opts *xhttp.RequestOptions) (*xhttp.Response, error) {
			// Before-logic
			fmt.Printf("Before: %s %s\n", req.Method, req.URL)

			// Call the next handler
			resp, err := next(req, opts)

			// After-logic
			fmt.Printf("After: %s %s\n", req.Method, req.URL)

			return resp, err
		}
	}
	resp, err := xhttp.Fetch(context.Background(), "GET", "https://github.com/", xhttp.WithMiddleware(logicMiddleware))

	a.Nil(err)
	a.Equal(resp.StatusCode, 200)
}

func TestShutdown(t *testing.T) {
	a := assert.New(t)
	logicMiddleware := func(next xhttp.HandlerFunc) xhttp.HandlerFunc {
		return func(req *xhttp.Request, opts *xhttp.RequestOptions) (*xhttp.Response, error) {
			// Before-logic
			fmt.Printf("Before: %s %s\n", req.Method, req.URL)

			// Call the next handler
			resp, err := next(req, opts)

			// After-logic
			fmt.Printf("After: %s %s %v\n", req.Method, req.URL, err)

			return resp, err
		}
	}
	_, err := xhttp.Fetch(context.Background(), "GET", "https://github.com/", xhttp.WithMiddleware(logicMiddleware))
	a.Nil(err)
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			_, err := xhttp.Fetch(context.Background(), "GET", "https://github.com/", xhttp.WithMiddleware(logicMiddleware))
			a.Equal(err, xhttp.ErrShutdown)
		}(i, &wg)
	}
	xhttp.Shutdown(context.Background())
	wg.Wait()
}
