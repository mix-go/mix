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

	DB := New(db)

	rows, _, err := DB.Query("SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}
	bar := rows[0].Get("bar").String()
	fmt.Println(bar)

	a.Equal(bar, "2022-04-14 23:49:48")
}

type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

func (t Test) TableName() string {
	return "xsql"
}

func TestInsert(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	DB := New(db)

	test := Test{
		Id:  0,
		Foo: "test",
		Bar: time.Now(),
	}
	_, l, err := DB.Insert(&test)
	fmt.Println(l)

	a.Empty(err)
}

func TestBatchInsert(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	DB := New(db)

	test := []Test{
		{
			Id:  0,
			Foo: "test",
			Bar: time.Now(),
		},
		{
			Id:  0,
			Foo: "test",
			Bar: time.Now(),
		},
	}
	_, l, err := DB.BatchInsert(&test)
	fmt.Println(l)

	a.Empty(err)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	DB := New(db)

	test := Test{
		Id:  999,
		Foo: "test update",
		Bar: time.Now(),
	}
	_, l, err := DB.Update(&test, "id = ?", 10)
	fmt.Println(l)

	a.Empty(err)
}

func TestFirst(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	DB := New(db)

	var test Test
	_, err = DB.First(&test, "SELECT * FROM xsql")
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

	DB := New(db)

	var tests []Test
	_, err = DB.Find(&tests, "SELECT * FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0800 CST} {Id:2 Foo:v1 Bar:2022-04-14 23:50:00 +0800 CST}]`)
}
