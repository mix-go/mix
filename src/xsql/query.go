package xsql

import (
	"time"
)

type query struct {
	Query
}

func (t *query) Fetch(query string, args []interface{}, opts *sqlOptions) (*Fetcher, error) {
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
		opts.doDebug(l)
		return nil, err
	}

	f := &Fetcher{
		r:       r,
		log:     l,
		options: opts,
	}
	return f, err
}
