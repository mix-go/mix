package xhttp

import (
	"github.com/avast/retry-go"
	"github.com/mix-go/xutil/xconv"
	"net/http"
	"time"
)

var DefaultOptions = newDefaultOptions()

func newDefaultOptions() requestOptions {
	return requestOptions{
		Header:  http.Header{},
		Timeout: time.Second * 5,
	}
}

type requestOptions struct {
	Header http.Header

	Body Body

	// 默认: time.Second * 5
	Timeout time.Duration

	DebugFunc DebugFunc

	// Retry
	RetryIfFunc  RetryIfFunc
	RetryOptions []retry.Option
}

func mergeOptions(opts []RequestOption) *requestOptions {
	cp := DefaultOptions // copy
	for _, o := range opts {
		o.apply(&cp)
	}
	return &cp
}

type RequestOption interface {
	apply(*requestOptions)
}

type funcRequestOption struct {
	f func(*requestOptions)
}

func (fdo *funcRequestOption) apply(do *requestOptions) {
	fdo.f(do)
}

func WithHeader(header http.Header) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Header = header
	}}
}

func WithContentType(contentType string) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Header.Set("Content-Type", contentType)
	}}
}

func WithTimeout(timeout time.Duration) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Timeout = timeout
	}}
}

func WithBody(body Body) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Body = body
	}}
}

func WithBodyBytes(body []byte) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Body = body
	}}
}

func WithBodyString(body string) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.Body = xconv.StringToBytes(body)
	}}
}

func WithDebugFunc(f DebugFunc) RequestOption {
	return &funcRequestOption{func(opts *requestOptions) {
		opts.DebugFunc = f
	}}
}

func WithRetry(f RetryIfFunc, opts ...retry.Option) RequestOption {
	return &funcRequestOption{func(o *requestOptions) {
		o.RetryIfFunc = f
		o.RetryOptions = opts
	}}
}
