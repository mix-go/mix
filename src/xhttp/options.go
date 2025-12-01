package xhttp

import (
	"github.com/avast/retry-go"
	"net/http"
	"time"
)

var DefaultOptions = newDefaultOptions()

func newDefaultOptions() RequestOptions {
	return RequestOptions{
		Header:  http.Header{},
		Timeout: time.Second * 5,
	}
}

type RequestOptions struct {
	Header http.Header

	Body Body

	// 默认: time.Second * 5
	Timeout time.Duration

	DebugFunc DebugFunc

	// Retry
	RetryIfFunc  RetryIfFunc
	RetryOptions []retry.Option

	Middlewares []Middleware
}

func mergeOptions(c *Client, opts []RequestOption) *RequestOptions {
	cp := c.DefaultOptions // copy
	for _, o := range opts {
		o.apply(&cp)
	}
	return &cp
}

type RequestOption interface {
	apply(*RequestOptions)
}

type funcRequestOption struct {
	f func(*RequestOptions)
}

func (fdo *funcRequestOption) apply(do *RequestOptions) {
	fdo.f(do)
}

func WithHeader(header http.Header) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Header = header
	}}
}

func WithContentType(contentType string) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Header.Set("Content-Type", contentType)
	}}
}

func WithTimeout(timeout time.Duration) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Timeout = timeout
	}}
}

func WithBody(body Body) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Body = body
	}}
}

func WithBodyBytes(body []byte) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Body = body
	}}
}

func WithBodyString(body string) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.Body = StringToBytes(body)
	}}
}

func WithDebugFunc(f DebugFunc) RequestOption {
	return &funcRequestOption{func(opts *RequestOptions) {
		opts.DebugFunc = f
	}}
}

func WithRetry(f RetryIfFunc, opts ...retry.Option) RequestOption {
	return &funcRequestOption{func(o *RequestOptions) {
		o.RetryIfFunc = f
		o.RetryOptions = opts
	}}
}

func WithMiddleware(middlewares ...Middleware) RequestOption {
	return &funcRequestOption{func(o *RequestOptions) {
		o.Middlewares = append(o.Middlewares, middlewares...)
	}}
}
