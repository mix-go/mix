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
	cli := http.Client{
		Timeout: opt.Timeout,
	}
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
	startTime := time.Now()
	r, err := cli.Do(req)
	if err != nil {
		debug(opt, time.Now().Sub(startTime), req, nil, err)
		return nil, err
	}
	resp := newResponse(r)
	debug(opt, time.Now().Sub(startTime), req, resp, nil)
	return resp, nil
}

func Do(req *http.Request, opts ...RequestOption) (*Response, error) {
	opt := getOptions(opts)
	cli := http.Client{
		Timeout: opt.Timeout,
	}
	startTime := time.Now()
	r, err := cli.Do(req)
	if err != nil {
		debug(opt, time.Now().Sub(startTime), req, nil, err)
		return nil, err
	}
	resp := newResponse(r)
	debug(opt, time.Now().Sub(startTime), req, resp, nil)
	return resp, nil
}

func debug(opt *requestOptions, duration time.Duration, req *http.Request, resp *Response, err error) {
	if opt.DebugFunc == nil {
		return
	}

	l := &Log{
		Duration: duration,
		Request:  req,
		Response: resp,
		Error:    err,
	}
	opt.DebugFunc(l)
}
