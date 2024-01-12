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
	var errorLog []error
	err = retry.Do(
		func() error {
			resp, err = f()
			if opts.RetryIfFunc != nil {
				err := opts.RetryIfFunc(resp, err)
				if err != nil {
					errorLog = append(errorLog, err)
					if errors.Is(err, ErrAbortRetry) {
						return nil
					}
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
	if len(errorLog) > 0 {
		return nil, errors.Join(errorLog...)
	}
	return resp, nil
}
