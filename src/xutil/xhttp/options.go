package xhttp

import (
	"net/http"
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
	Body   Body

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
