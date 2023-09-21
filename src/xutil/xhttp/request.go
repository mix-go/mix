package xhttp

import (
	"github.com/mix-go/xutil/xconv"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	*http.Response
	Body Body
}

type Body []byte

func (t Body) String() string {
	return xconv.BytesToString(t)
}

func newResponse(r *http.Response) *Response {
	if r == nil {
		return nil
	}
	resp := &Response{
		Response: r,
	}
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return resp
	}
	resp.Body = b
	return resp
}

func Request(method string, u string, opts ...RequestOption) (*Response, error) {
	opt := getOptions(opts)
	URL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: method,
		URL:    URL,
		Header: opt.Header,
	}
	if opt.Body != nil {
		req.Body = io.NopCloser(strings.NewReader(opt.Body.String()))
	}

	if opt.RetryOptions != nil {
		return doRetry(opt, func() (*Response, error) {
			return do(opt, req)
		})
	}
	return do(opt, req)
}

func Do(req *http.Request, opts ...RequestOption) (*Response, error) {
	opt := getOptions(opts)

	if opt.RetryOptions != nil {
		return doRetry(opt, func() (*Response, error) {
			return do(opt, req)
		})
	}
	return do(opt, req)
}

func do(opt *requestOptions, req *http.Request) (*Response, error) {
	cli := http.Client{
		Timeout: opt.Timeout,
	}
	startTime := time.Now()
	r, err := cli.Do(req)
	if err != nil {
		doDebug(opt, time.Now().Sub(startTime), req, nil, err)
		return nil, err
	}
	resp := newResponse(r)
	doDebug(opt, time.Now().Sub(startTime), req, resp, nil)
	return resp, nil
}
