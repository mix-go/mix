package xsql

import (
	"database/sql"
	"time"
)

var TimeParselayout = "2006-01-02 15:04:05"

func Query(db *sql.DB, query string, args ...interface{}) ([]Row, *Log, error) {
	startTime := time.Now()
	r, err := db.Query(query, args...)
	l := &Log{
		SQL:  query,
		Args: args,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}
	f := Fetcher{R: r}
	rows, err := f.Rows()
	return rows, l, err
}

func Find(r *sql.Rows, i interface{}) error {
	f := &Fetcher{R: r}
	return f.Find(i)
}

func First(r *sql.Rows, i interface{}) error {
	f := &Fetcher{R: r}
	return f.First(i)
}

func Insert(db *sql.DB, table string, data interface{}, opts ...Option) (sql.Result, *Log, error) {
	e := &Executor{
		DB: db,
	}
	return e.Insert(table, data, opts...)
}

func BatchInsert(db *sql.DB, table string, data interface{}, opts ...Option) (sql.Result, *Log, error) {
	e := &Executor{
		DB: db,
	}
	return e.BatchInsert(table, data, opts...)
}
