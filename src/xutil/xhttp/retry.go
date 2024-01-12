package xhttp

import (
	"errors"
	"github.com/avast/retry-go"
)

var ErrAbortRetry = errors.New("xhttp: abort further retries")

type RetryIfFunc func(*XResponse, error) error

func doRetry(opts *requestOptions, f func() (*XResponse, error)) (*XResponse, error) {
	var resp *XResponse
	var err error
	var lastErr error
	err = retry.Do(
		func() error {
			resp, err = f()
			if opts.RetryIfFunc != nil {
				err := opts.RetryIfFunc(resp, err)
				if err != nil && errors.Is(err, ErrAbortRetry) {
					lastErr = err
					return nil
				}
				return err
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
	if lastErr != nil {
		return nil, lastErr
	}
	return resp, nil
}
