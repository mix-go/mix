package xhttp

import (
	"time"
)

type Log struct {
	Duration time.Duration `json:"duration"`
	Request  *XRequest     `json:"request"`  // The XRequest.RetryAttempts field records the number of retries that have been completed.
	Response *XResponse    `json:"response"` // If request error this field is equal to nil
	Error    error         `json:"error"`
}

type DebugFunc func(l *Log)

func doDebug(opts *requestOptions, duration time.Duration, req *XRequest, resp *XResponse, err error) {
	if opts.DebugFunc == nil {
		return
	}

	l := &Log{
		Duration: duration,
		Request:  req,
		Response: resp,
		Error:    err,
	}
	opts.DebugFunc(l)
}
