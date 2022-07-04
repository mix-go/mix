package xsql

import "database/sql"

type Tx struct {
	raw *sql.Tx
	*DB
}

func (t *Tx) Commit() error {
	return t.raw.Commit()
}

func (t *Tx) Rollback() error {
	return t.raw.Rollback()
}
