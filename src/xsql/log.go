package xsql

import "time"

type Log struct {
	Time     time.Duration `json:"time"`
	SQL      string        `json:"sql"`
	Bindings []interface{} `json:"bindings"`
}

type DebugFunc func(l *Log)
