package xsql

import (
	"database/sql"
)

var Tag = "xsql"

var DefaultTimeLayout = "2006-01-02 15:04:05"

// DefaultTimeFunc
// mysql: return placeholder
// oracle: return fmt.Sprintf("TO_TIMESTAMP(%s, 'SYYYY-MM-DD HH24:MI:SS:FF6')", placeholder)
var DefaultTimeFunc = func(placeholder string) string {
	return placeholder
}

type TimeFunc func(placeholder string) string

type DB struct {
	Options  Options
	raw      *sql.DB
	executor executor
	query    query
}

// New
// opts 以最后一个为准
func New(db *sql.DB, opts ...Options) *DB {
	o := Options{}
	for _, v := range opts {
		o = v
	}
	return &DB{
		Options: o,
		raw:     db,
		executor: executor{
			Executor: db,
		},
		query: query{
			Query: db,
		},
	}
}

func (t *DB) Insert(data interface{}, opts ...Options) (sql.Result, error) {
	for _, o := range opts {
		t.Options.InsertKey = o.InsertKey
	}
	return t.executor.Insert(data, &t.Options)
}

func (t *DB) BatchInsert(data interface{}, opts ...Options) (sql.Result, error) {
	for _, o := range opts {
		t.Options.InsertKey = o.InsertKey
	}
	return t.executor.BatchInsert(data, &t.Options)
}

func (t *DB) Update(data interface{}, expr string, args ...interface{}) (sql.Result, error) {
	return t.executor.Update(data, expr, args, &t.Options)
}

func (t *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.executor.Exec(query, args, &t.Options)
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
	f, err := t.query.Fetch(query, args, &t.Options)
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
	f, err := t.query.Fetch(query, args, &t.Options)
	if err != nil {
		return err
	}
	if err := f.Find(i); err != nil {
		return err
	}
	return nil
}

func (t *DB) First(i interface{}, query string, args ...interface{}) error {
	f, err := t.query.Fetch(query, args, &t.Options)
	if err != nil {
		return err
	}
	if err := f.First(i); err != nil {
		return err
	}
	return nil
}
