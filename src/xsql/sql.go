package xsql

import (
	"database/sql"
)

var DefaultTimeLayout = "2006-01-02 15:04:05"

type Database struct {
	Options
	executor
	query
}

// New
// opts 以最后一个为准
func New(db *sql.DB, opts ...Options) *Database {
	o := Options{}
	for _, v := range opts {
		o = v
	}
	return &Database{
		Options: o,
		executor: executor{
			DB: db,
		},
		query: query{
			DB: db,
		},
	}
}

func (t *Database) Insert(data interface{}, opts ...Options) (sql.Result, error) {
	for _, o := range opts {
		t.Options.InsertKey = o.InsertKey
	}
	return t.executor.Insert(data, &t.Options)
}

func (t *Database) BatchInsert(data interface{}, opts ...Options) (sql.Result, error) {
	for _, o := range opts {
		t.Options.InsertKey = o.InsertKey
	}
	return t.executor.BatchInsert(data, &t.Options)
}

func (t *Database) Update(data interface{}, expr string, args ...interface{}) (sql.Result, error) {
	return t.executor.Update(data, expr, args, &t.Options)
}

func (t *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.executor.Exec(query, args, &t.Options)
}

func (t *Database) Query(query string, args ...interface{}) ([]Row, error) {
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

func (t *Database) Find(i interface{}, query string, args ...interface{}) error {
	f, err := t.query.Fetch(query, args, &t.Options)
	if err != nil {
		return err
	}
	if err := f.Find(i); err != nil {
		return err
	}
	return nil
}

func (t *Database) First(i interface{}, query string, args ...interface{}) error {
	f, err := t.query.Fetch(query, args, &t.Options)
	if err != nil {
		return err
	}
	if err := f.First(i); err != nil {
		return err
	}
	return nil
}
