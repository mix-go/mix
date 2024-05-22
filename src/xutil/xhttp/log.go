package xhttp

import (
	"time"
)

type Log struct {
	Duration time.Duration `json:"duration"`
	Request  *Request      `json:"request"`  // The Request.RetryAttempts field records the number of retry attempts
	Response *Response     `json:"response"` // If request error this field is equal to nil
	Error    error         `json:"error"`
}

type DebugFunc func(l *Log)

func (t *Client) doDebug(opts *RequestOptions, duration time.Duration, req *Request, resp *Response, err error) {
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
