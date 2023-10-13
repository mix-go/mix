package xhttp_test

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mix-go/xutil/xhttp"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.Request("GET", url)

	a.Equal(resp.StatusCode, 200)
	a.Nil(err)
}

func TestRequestPOST(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := xhttp.Request("POST", url, xhttp.WithBodyString("abc"), xhttp.WithContentType("application/json"))

	a.Equal(resp.StatusCode, 404)
	a.Nil(err)
}

func TestRequestError(t *testing.T) {
	a := assert.New(t)

	url := "https://aaaaa.com/"
	resp, err := xhttp.Request("GET", url)

	a.Nil(resp)
	a.NotNil(err)
}

func TestDebugAndRetry(t *testing.T) {
	a := assert.New(t)

	xhttp.DefaultOptions.DebugFunc = func(l *xhttp.Log) {
		log.Printf("%+v %+v %+v %+v\n", l.Duration, l.Request, l.Response, l.Error)
	}

	url := "https://aaaaa.com/"
	retryIf := func(resp *xhttp.XResponse, err error) error {
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status_code: %d", resp.StatusCode)
		}
		return nil
	}
	resp, err := xhttp.Request("GET", url, xhttp.WithRetry(retryIf, retry.Attempts(2)))

	a.Nil(resp)
	a.NotNil(err)
	a.Contains(err.Error(), "attempts fail")
}
