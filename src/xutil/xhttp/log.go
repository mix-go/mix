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
