package di

import (
	"database/sql"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xdi"
	"github.com/mix-go/xsql"
)

func init() {
	obj := xdi.Object{
		Name: "xsql",
		New: func() (i interface{}, e error) {
			db, err := sql.Open("mysql", dotenv.Getenv("DATABASE_DSN").String())
			if err != nil {
				return nil, err
			}
			return xsql.New(db), nil
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func Xsql() (db *xsql.DB) {
	if err := xdi.Populate("xsql", &db); err != nil {
		panic(err)
	}
	return
}
