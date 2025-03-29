package xsql_test

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mix-go/xsql"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"testing"
	"time"
)

type Enum int32

type Test struct {
	Id   int       `xsql:"id"`
	Foo  string    `xsql:"foo"`
	Bar  time.Time `xsql:"bar"`
	Bool bool      `xsql:"bool" json:"-"`
	Enum Enum      `xsql:"enum" json:"-"`
}

type TestJsonStruct struct {
	Test
	Json JsonItem `xsql:"json"`
}

type TestJsonStructPtr struct {
	Test
	Json *JsonItem `xsql:"json"`
}

type TestJsonSlice struct {
	Test
	Json []int `xsql:"json"`
}

type JsonItem struct {
	Foo string `xsql:"foo"`
}

func (t *Test) TableName() string {
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

type Test3 struct {
	Id  int                    `xsql:"id"`
	Foo string                 `xsql:"foo"`
	Bar *timestamppb.Timestamp `xsql:"bar"`
}

func (t Test3) TableName() string {
	return "xsql"
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
	DB := newDB()

	b, err := os.ReadFile("./xsql.sql")
	a.Nil(err)
	_, err = DB.Exec(string(b))
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
	log.Println(bar)

	a.Equal(bar, "2022-04-12 23:50:00")
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
			Foo: "test_update",
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
		Foo: "test_update_1",
		Bar: time.Now(),
	}
	_, err := DB.Update(&test, "id = ?", 999)

	a.Empty(err)
}

func TestUpdateColumns(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	data := map[string]interface{}{
		"foo": "test_update_2",
	}
	_, err := DB.Model(&Test{}).Update(data, "id = ?", 8)

	a.Empty(err)
}

func TestUpdateTagValuesMap(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{}
	data, err := xsql.TagValuesMap(DB.Options.Tag, &test,
		xsql.TagValues{
			{&test.Foo, "test_update_3"},
		},
	)
	a.Empty(err)

	_, err = DB.Model(&test).Update(data, "id = ?", 8)
	a.Empty(err)
}

func TestEmbeddingUpdateTagValuesMap(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := EmbeddingTest{}
	data, err := xsql.TagValuesMap(DB.Options.Tag, &test,
		xsql.TagValues{
			{&test.Foo, "test_update_4"},
		},
	)
	a.Empty(err)

	_, err = DB.Model(&test).Update(data, "id = ?", 8)
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
	err := DB.First(&test, "SELECT * FROM ${TABLE}")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
	// bool
	a.Equal(test.Bool, true)
	// enum
	a.IsType(Enum(0), test.Enum)
	a.Equal(Enum(1), test.Enum)
}

func TestFirstPtr(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test *Test
	err := DB.First(&test, "SELECT * FROM ${TABLE}")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
	// bool
	a.Equal(test.Bool, true)
	// enum
	a.IsType(Enum(0), test.Enum)
	a.Equal(Enum(1), test.Enum)
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
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
}

func TestFirstPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(&test, "SELECT foo FROM xsql")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(test)
	a.Equal(string(b), "{\"Id\":0,\"Foo\":\"v\",\"Bar\":\"0001-01-01T00:00:00Z\"}")
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
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
}

func TestFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT * FROM ${TABLE} LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
}

func TestFindPtr(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []*Test
	err := DB.Find(&tests, "SELECT * FROM ${TABLE} LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
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
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
}

func TestFindPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT foo FROM xsql LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(tests)
	a.Equal(string(b), "[{\"Id\":0,\"Foo\":\"v\",\"Bar\":\"0001-01-01T00:00:00Z\"},{\"Id\":0,\"Foo\":\"v1\",\"Bar\":\"0001-01-01T00:00:00Z\"}]")
}

func TestFindTableKey(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(&tests, "SELECT * FROM ${TABLE} LIMIT 2")
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
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

func TestPbTimestamp(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	// Insert
	now := timestamppb.Now()
	log.Println(now.AsTime().Format(time.RFC3339))
	test := Test3{
		Id:  0,
		Foo: "test_pb_timestamp",
		Bar: now,
	}
	res, err := DB.Insert(&test)
	a.Empty(err)
	insertId, _ := res.LastInsertId()

	// First
	var test2 Test3
	err = DB.First(&test2, "SELECT * FROM xsql WHERE id = ?", insertId)
	if err != nil {
		log.Fatal(err)
	}
	// Timestamp
	a.IsType(&timestamppb.Timestamp{}, test2.Bar)
	a.Equal(test2.Bar.Seconds, now.Seconds)
}

func TestFetchPbJson(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	var test1 TestJsonStruct
	err := DB.First(&test1, "SELECT * FROM xsql WHERE id = 1")
	if err != nil {
		log.Fatal(err)
	}
	a.NotEmpty(test1.Json)

	var test2 TestJsonStructPtr
	err = DB.First(&test2, "SELECT * FROM xsql WHERE id = 1")
	if err != nil {
		log.Fatal(err)
	}
	a.NotEmpty(test2.Json)

	var test3 TestJsonSlice
	err = DB.First(&test3, "SELECT * FROM xsql WHERE id = 2")
	if err != nil {
		log.Fatal(err)
	}
	a.NotEmpty(test3.Json)
}

func TestInsertPbJson(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	test1 := TestJsonStruct{
		Test: Test{
			Id:   0,
			Foo:  "",
			Bar:  time.Time{},
			Bool: false,
			Enum: 0,
		},
		Json: JsonItem{Foo: `bar`},
	}
	_, err := DB.Insert(&test1)
	if err != nil {
		log.Fatal(err)
	}
	a.Empty(err)

	test2 := TestJsonStructPtr{
		Test: Test{
			Id:   0,
			Foo:  "",
			Bar:  time.Time{},
			Bool: false,
			Enum: 0,
		},
		Json: &JsonItem{Foo: `bar`},
	}
	_, err = DB.Insert(&test2)
	if err != nil {
		log.Fatal(err)
	}
	a.Empty(err)

	test3 := TestJsonSlice{
		Test: Test{
			Id:   0,
			Foo:  "",
			Bar:  time.Time{},
			Bool: false,
			Enum: 0,
		},
		Json: []int{1, 2, 3},
	}
	_, err = DB.Insert(&test3)
	if err != nil {
		log.Fatal(err)
	}
	a.Empty(err)
}
