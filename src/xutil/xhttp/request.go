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

func Request(method string, u string, body string, opts ...Options) ([]byte, error) {
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
		Body:   io.NopCloser(strings.NewReader(body)),
		Header: opt.Header,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return b, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return b, nil
}

func RequestString(method string, u string, body string, opts ...Options) (string, error) {
	b, err := Request(method, u, body, opts...)
	return string(b), err
}
