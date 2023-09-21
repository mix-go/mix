package xsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

func newDB() *DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&multiStatements=true")
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

func TestCreateTable(t *testing.T) {
	a := assert.New(t)
	q := `DROP TABLE IF EXISTS #xsql#;
CREATE TABLE #xsql# (
  #id# int unsigned NOT NULL AUTO_INCREMENT,
  #foo# varchar(255) DEFAULT NULL,
  #bar# datetime DEFAULT NULL,
  PRIMARY KEY (#id#)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
INSERT INTO #xsql# (#id#, #foo#, #bar#) VALUES (1, 'v', '2022-04-14 23:49:48');
INSERT INTO #xsql# (#id#, #foo#, #bar#) VALUES (2, 'v1', '2022-04-14 23:50:00');
`
	DB := newDB()
	_, err := DB.Exec(strings.ReplaceAll(q, "#", "`"))
	a.Empty(err)
}

func TestClear(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	_, err := DB.Exec("DELETE FROM xsql WHERE id > 2")

	a.Empty(err)
}

func TestDebugFunc(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  0,
		Foo: "test update",
		Bar: time.Now(),
	}
	_, err := DB.Update(&test, "id = ?", 0)

	a.Empty(err)
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

func TestFirst(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(&test, "SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48+08:00"}`)
}

type EmbeddingTest struct {
	Test
}

func TestEmbeddingFirst(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test EmbeddingTest
	err := DB.First(&test, "SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48+08:00"}`)
}

func TestFirstPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(&test, "SELECT foo FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", test), "{Id:0 Foo:v Bar:0001-01-01 00:00:00 +0000 UTC}")
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

func TestFindPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT foo FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:0 Foo:v Bar:0001-01-01 00:00:00 +0000 UTC} {Id:0 Foo:v1 Bar:0001-01-01 00:00:00 +0000 UTC}]`)
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
