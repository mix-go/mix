package xhttp_test

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mix-go/xutil/xhttp"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestNewRequest(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.NewRequest("GET", url)

	a.Equal(resp.StatusCode, 200)
	a.Nil(err)
}

func TestRequestPOST(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.NewRequest("POST", url, xhttp.WithBodyString("abc"), xhttp.WithContentType("application/json"))

	a.Equal(resp.StatusCode, 404)
	a.Nil(err)
}

func TestRequestError(t *testing.T) {
	a := assert.New(t)

	url := "https://aaaaa.com/"
	resp, err := xhttp.NewRequest("GET", url)

	a.Nil(resp)
	a.NotNil(err)
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
	resp, err := xhttp.NewRequest("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(2)))

	a.Nil(resp)
	a.NotNil(err)
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
	_, err := xhttp.NewRequest("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(3)))

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
	resp, err := xhttp.NewRequest("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(3)))

	a.Nil(resp)
	a.NotNil(err)
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
	resp, err := xhttp.NewRequest("GET", "https://github.com/", xhttp.WithMiddlewares(logicMiddleware))

	a.Equal(resp.StatusCode, 200)
	a.Nil(err)
}
