package xhttp

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
)

type RetryIfFunc func(*Response, error) error

type Error []error

func (t Error) Error() string {
	var logWithNumber []error
	for i, err := range t {
		logWithNumber = append(logWithNumber, fmt.Errorf("#%d: %s", i+1, err))
	}
	return errors.Join(logWithNumber...).Error()
}

func (t Error) HasAbortRetry() bool {
	for _, err := range t {
		if errors.Is(err, ErrAbortRetry) {
			return true
		}
	}
	return false
}

func doRetry(opts *RequestOptions, f func() (*Response, error)) (*Response, error) {
	var resp *Response
	var err error
	var errorLog Error
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
	if errorLog.HasAbortRetry() {
		return nil, &errorLog
	}
	return resp, nil
}
