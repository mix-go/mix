package xhttp

import (
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
}

func getOptions(opts []RequestOption) requestOptions {
	opt := DefaultOptions
	for _, o := range opts {
		o.apply(&opt)
	}
	return opt
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
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Header = header
	}}
}

func WithContentType(contentType string) RequestOption {
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Header.Set("Content-Type", contentType)
	}}
}

func WithTimeout(timeout time.Duration) RequestOption {
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Timeout = timeout
	}}
}

func WithBody(body Body) RequestOption {
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Body = body
	}}
}

func WithBodyBytes(body []byte) RequestOption {
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Body = body
	}}
}

func WithBodyString(body string) RequestOption {
	return &funcRequestOption{func(opt *requestOptions) {
		opt.Body = xconv.StringToBytes(body)
	}}
}
