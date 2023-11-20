package xsql_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mix-go/xsql"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

func (t Test) TableName() string {
	return "xsql"
}

type Test1 struct {
	Id int `xsql:"id"`
}

func (t Test1) TableName() string {
	return "xsql"
}

type Test2 struct {
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

type EmbeddingTest struct {
	Test1
	Test2
}

func newDB() *xsql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true&loc=UTC&multiStatements=true")
	if err != nil {
		log.Fatal(err)
	}
	return xsql.New(
		db,
		xsql.WithDebugFunc(func(l *xsql.Log) {
			log.Println(l)
		}),
		xsql.WithTimeLocation(time.UTC),
	)
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
		Foo: "test",
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

func TestEmbeddingInsert(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := EmbeddingTest{
		Test1: Test1{
			Id: 0,
		},
		Test2: Test2{
			Foo: "test",
			Bar: time.Now(),
		},
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
			Foo: "test1",
			Bar: time.Now(),
		},
		{
			Id:  0,
			Foo: "test2",
			Bar: time.Now(),
		},
	}
	_, err := DB.BatchInsert(&tests)

	a.Empty(err)
}

func TestEmbeddingBatchInsert(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	tests := []EmbeddingTest{
		{
			Test1: Test1{
				Id: 0,
			},
			Test2: Test2{
				Foo: "test1",
				Bar: time.Now(),
			},
		},
		{
			Test1: Test1{
				Id: 0,
			},
			Test2: Test2{
				Foo: "test2",
				Bar: time.Now(),
			},
		},
	}
	_, err := DB.BatchInsert(&tests)

	a.Empty(err)
}

func TestEmbeddingUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newDB()
	test := EmbeddingTest{
		Test1: Test1{
			Id: 999,
		},
		Test2: Test2{
			Foo: "test update",
			Bar: time.Now(),
		},
	}
	_, err := DB.Update(&test, "id = ?", 8)

	a.Empty(err)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  8,
		Foo: "test update 1",
		Bar: time.Now(),
	}
	_, err := DB.Update(&test, "id = ?", 999)

	a.Empty(err)
}

func TestUpdateColumns(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := map[string]interface{}{
		"foo": "test update 2",
	}
	_, err := DB.Model(&Test{}).Update(test, "id = ?", 8)

	a.Empty(err)
}

func TestDelete(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  8,
		Foo: "test",
		Bar: time.Now(),
	}
	_, err := DB.Model(&test).Delete("id = ?", test.Id)
	a.Empty(err)

	_, err = DB.Model(&Test{}).Delete("id = ?", 8)
	a.Empty(err)
}

func TestExec(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	_, err := DB.Exec("DELETE FROM xsql WHERE id = ?", 7)

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
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48Z"}`)
}

func TestFirstEmbedding(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test EmbeddingTest
	err := DB.First(&test, "SELECT * FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48Z"}`)
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

func TestFirstTableKey(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(&test, "SELECT * FROM ${TABLE}")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48Z"}`)
}

func TestFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT * FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0000 UTC} {Id:2 Foo:v1 Bar:2022-04-14 23:50:00 +0000 UTC}]`)
}

func TestEmbeddingFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []EmbeddingTest
	err := DB.Find(&tests, "SELECT * FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-14T23:49:48Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-14T23:50:00Z"}]`)
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

func TestFindTableKey(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT * FROM ${TABLE} LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:2022-04-14 23:49:48 +0000 UTC} {Id:2 Foo:v1 Bar:2022-04-14 23:50:00 +0000 UTC}]`)
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
