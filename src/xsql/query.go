package xsql

import (
	"time"
)

type query struct {
	Query
}

func (t *query) Fetch(query string, args []interface{}, opts *Options) (*Fetcher, error) {
	startTime := time.Now()
	r, err := t.Query.Query(query, args...)
	l := &Log{
		Duration:     time.Now().Sub(startTime),
		SQL:          query,
		Bindings:     args,
		RowsAffected: 0,
		Error:        err,
	}
	if err != nil {
		if opts.DebugFunc != nil {
			opts.DebugFunc(l)
		}
		return nil, err
	}

	f := &Fetcher{
		R:       r,
		Log:     l,
		Options: opts,
	}
	return f, err
}
