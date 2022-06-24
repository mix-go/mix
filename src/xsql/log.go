package xsql

import "time"

type Log struct {
	SQL  string
	Args []interface{}
	Time time.Duration
}
