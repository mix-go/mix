package xsql

import (
	"database/sql"
)

type TimeFunc func(placeholder string) string

type DB struct {
	Options  *sqlOptions
	raw      *sql.DB
	executor executor
	query    query
}

func New(db *sql.DB, opts ...SqlOption) *DB {
	return &DB{
		Options: mergeOptions(opts),
		raw:     db,
		executor: executor{
			Executor: db,
		},
		query: query{
			Query: db,
		},
	}
}

func (t *DB) mergeOptions(opts []SqlOption) *sqlOptions {
	opt := *t.Options // copy
	for _, o := range opts {
		o.apply(&opt)
	}
	return &opt
}

func (t *DB) Insert(data interface{}, opts ...SqlOption) (sql.Result, error) {
	return t.executor.Insert(data, t.mergeOptions(opts))
}

func (t *DB) BatchInsert(data interface{}, opts ...SqlOption) (sql.Result, error) {
	return t.executor.BatchInsert(data, t.mergeOptions(opts))
}

func (t *DB) Update(data interface{}, expr string, args ...interface{}) (sql.Result, error) {
	return t.executor.Update(data, expr, args, t.Options)
}

func (t *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.executor.Exec(query, args, t.Options)
}

func (t *DB) Begin() (*Tx, error) {
	tx, err := t.raw.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{
		raw: tx,
		DB: &DB{
			Options: t.Options,
			executor: executor{
				Executor: tx,
			},
			query: query{
				Query: tx,
			},
		},
	}, nil
}

func (t *DB) Query(query string, args ...interface{}) ([]Row, error) {
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return nil, err
	}
	r, err := f.Rows()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t *DB) Find(i interface{}, query string, args ...interface{}) error {
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return err
	}
	if err := f.Find(i); err != nil {
		return err
	}
	return nil
}

func (t *DB) First(i interface{}, query string, args ...interface{}) error {
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return err
	}
	if err := f.First(i); err != nil {
		return err
	}
	return nil
}
