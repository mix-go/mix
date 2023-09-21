package xsql

import "time"

type Log struct {
	Duration     time.Duration `json:"time"`
	SQL          string        `json:"sql"`
	Bindings     []interface{} `json:"bindings"`
	RowsAffected int64         `json:"rowsAffected"`
	Error        error         `json:"error"`
}

type DebugFunc func(l *Log)
