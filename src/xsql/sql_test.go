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

func newDB() *DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	opts := Options{
		DebugFunc: func(l *Log) {
			log.Println(l)
		},
	}
	return New(db, opts)
}

func TestQuery(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	rows, err := DB.Query("SELECT * FROM xsql")
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

	DB := newDB()

	test := Test{
		Id:  0,
		Foo: "test",
		Bar: time.Now(),
	}
	_, err := DB.Insert(&test)

	a.Empty(err)
}

func TestBatchInsert(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	tests := []Test{
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
	_, err := DB.BatchInsert(&tests)

	a.Empty(err)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  999,
		Foo: "test update",
		Bar: time.Now(),
	}
	_, err := DB.Update(&test, "id = ?", 10)

	a.Empty(err)
}

func TestExec(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	_, err := DB.Exec("DELETE FROM xsql WHERE id = ?", 10)

	a.Empty(err)
}

func TestDebugFunc(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  999,
		Foo: "test update",
		Bar: time.Now(),
	}
	_, err := DB.Update(&test, "id = ?", 10)

	a.Empty(err)
}

func TestFirst(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(&test, "SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", test), "{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0800 CST}")
}

func TestFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT * FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0800 CST} {Id:2 Foo:v1 Bar:2022-04-14 23:50:00 +0800 CST}]`)
}

func TestTxCommit(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	tx, _ := DB.Begin()

	test := Test{
		Id:  0,
		Foo: "test",
		Bar: time.Now(),
	}
	_, err := tx.Insert(&test)
	a.Empty(err)

	err = tx.Commit()
	a.Empty(err)
}

func TestTxRollback(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	tx, _ := DB.Begin()

	test := Test{
		Id:  0,
		Foo: "test",
		Bar: time.Now(),
	}
	_, err := tx.Insert(&test)
	a.Empty(err)

	err = tx.Rollback()
	a.Empty(err)
}
