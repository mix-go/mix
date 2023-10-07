package xsql

import "database/sql"

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Query interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Table interface {
	TableName() string
}
