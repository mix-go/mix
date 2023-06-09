package xhttp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var DefaultOptions = newDefaultOptions()

func newDefaultOptions() Options {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	return Options{
		Header:  h,
		Timeout: time.Second * 5,
	}
}

type Options struct {
	Header http.Header
	Body   string

	// 默认: time.Second * 5
	Timeout time.Duration
}

func newOptions(opts []Options) Options {
	current := DefaultOptions
	for _, v := range opts {
		current = v
	}
	if current.Timeout == 0 {
		current.Timeout = DefaultOptions.Timeout
	}
	return current
}

type Response struct {
	*http.Response
	Body string
}

func newResponse(r *http.Response) *Response {
	resp := &Response{
		Response: r,
	}
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return resp
	}
	resp.Body = string(b)
	return resp
}

func Request(method string, u string, opts ...Options) (*Response, error) {
	opt := newOptions(opts)
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
		Body:   io.NopCloser(strings.NewReader(opt.Body)),
		Header: opt.Header,
	}
	r, err := cli.Do(req)
	resp := newResponse(r)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return resp, nil
}
