package xhttp

import (
	"time"
)

type Log struct {
	Duration time.Duration `json:"duration"`
	Request  *XRequest     `json:"request"`
	Response *XResponse    `json:"response"`
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
