package xhttp

import (
	"github.com/avast/retry-go"
)

type RetryIfFunc func(*Response, error) error

func doRetry(opt *requestOptions, f func() (*Response, error)) (*Response, error) {
	var resp *Response
	var err error
	err = retry.Do(
		func() error {
			resp, err = f()
			if opt.RetryIfFunc != nil {
				return opt.RetryIfFunc(resp, err)
			}
			if err != nil {
				return err
			}
			return nil
		},
		opt.RetryOptions...,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
