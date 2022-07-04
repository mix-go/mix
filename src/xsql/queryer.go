package xsql

import (
	"time"
)

type query struct {
	Query
}

func (t *query) Fetch(query string, args []interface{}, opts *Options) (*Fetcher, error) {
	var debugFunc DebugFunc
	if opts.DebugFunc != nil {
		debugFunc = opts.DebugFunc
	}
	startTime := time.Now()
	r, err := t.Query.Query(query, args...)
	l := &Log{
		SQL:      query,
		Bindings: args,
		Time:     time.Now().Sub(startTime),
	}
	if debugFunc != nil {
		debugFunc(l)
	}
	if err != nil {
		return nil, err
	}
	f := &Fetcher{
		R:       r,
		Options: opts,
	}
	return f, err
}
