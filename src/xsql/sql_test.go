package xsql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := Query(db, "SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}
	bar := rows[0].Get("bar").String()

	a.Equal(bar, "2022-04-14 23:49:48")
}

type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

func TestFirst(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	var test Test
	err = First(rows, &test)
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", test), "{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0800 CST}")
}

func TestFind(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}
	var tests []Test
	err = Find(rows, &tests)
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0800 CST} {Id:2 Foo:v1 Bar:2022-04-14 23:50:00 +0800 CST}]`)
}
