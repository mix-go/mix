package xsql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
	"time"
)

type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar string    `xsql:"bar"`
	Orz time.Time `xsql:"orz,2006-01-02 15:04:05"`
}

func TestFetch(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		log.Fatal(err)
	}

	f := Fetch(rows)
	var a Test
	err = f.First(&a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", a)
}
