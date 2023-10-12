package xhttp

import (
	"github.com/avast/retry-go"
)

type RetryIfFunc func(*XResponse, error) error

func doRetry(opts *requestOptions, f func() (*XResponse, error)) (*XResponse, error) {
	var resp *XResponse
	var err error
	err = retry.Do(
		func() error {
			resp, err = f()
			if opts.RetryIfFunc != nil {
				return opts.RetryIfFunc(resp, err)
			}
			if err != nil {
				return err
			}
			return nil
		},
		opts.RetryOptions...,
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
