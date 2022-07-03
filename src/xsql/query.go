package xsql

import (
	"database/sql"
	"time"
)

type queryManager struct {
	DB *sql.DB
}

func (t *queryManager) Fetch(query string, args []interface{}, opts *Options) (*Fetcher, *Log, error) {
	startTime := time.Now()
	r, err := t.DB.Query(query, args...)
	l := &Log{
		SQL:  query,
		Args: args,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}
	f := &Fetcher{
		R:       r,
		Options: opts,
	}
	return f, l, err
}
