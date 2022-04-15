package xsql

import (
	"database/sql"
)

var TimeParselayout = "2006-01-02 15:04:05"

func Query(db *sql.DB, query string, args ...interface{}) ([]Row, error) {
	r, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	f := Fetcher{r: r}
	return f.Rows()
}

func Find(r *sql.Rows, i interface{}) error {
	f := &Fetcher{r: r}
	return f.Find(i)
}

func First(r *sql.Rows, i interface{}) error {
	f := &Fetcher{r: r}
	return f.First(i)
}
