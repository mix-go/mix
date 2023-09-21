package xhttp

import (
	"net/http"
	"time"
)

type Log struct {
	Duration time.Duration `json:"duration"`
	Request  *http.Request `json:"request"`
	Response *Response     `json:"response"`
	Error    error         `json:"error"`
}

type DebugFunc func(l *Log)

func doDebug(opt *requestOptions, duration time.Duration, req *http.Request, resp *Response, err error) {
	if opt.DebugFunc == nil {
		return
	}

	l := &Log{
		Duration: duration,
		Request:  req,
		Response: resp,
		Error:    err,
	}
	opt.DebugFunc(l)
}
