package xsql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mix-go/xsql"
	"github.com/mix-go/xsql/testdata"
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

type Test1 struct {
	Id int `xsql:"id"`
}

type Test2 struct {
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

type Test3 struct {
	Id  int                    `xsql:"id"`
	Foo string                 `xsql:"foo"`
	Bar *timestamppb.Timestamp `xsql:"bar"`
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

type TestJsonSlicePtr struct {
	Test
	Json []*JsonItem `xsql:"json"`
}

type JsonItem struct {
	Foo string `xsql:"foo"`
}

type TestEmbedding struct {
	Test1
	Test2
}

type TestPbStruct struct {
	testdata.Device
}

func (t *Test) TableName() string {
	return "xsql"
}

func (t *Test1) TableName() string {
	return "xsql"
}

func (t *Test2) TableName() string {
	return "xsql"
}

func (t *Test3) TableName() string {
	return "xsql"
}

func (t *TestEmbedding) TableName() string {
	return "xsql"
}

func (t *TestJsonStruct) TableName() string {
	return "xsql"
}

func (t *TestJsonStructPtr) TableName() string {
	return "xsql"
}

func (t *TestJsonSlice) TableName() string {
	return "xsql"
}

func (t *TestPbStruct) TableName() string {
	return "devices"
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
	err = DB.Exec(context.Background(), string(b)).Error
	a.Nil(err)
}

func TestClear(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	err := DB.Exec(context.Background(), "DELETE FROM xsql WHERE id > 2").Error
	a.Nil(err)
}

func TestDebugFunc(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  0,
		Foo: "test",
		Bar: time.Now(),
	}
	err := DB.Update(context.Background(), &test, "id = ?", 0).Error
	a.Nil(err)
}

func TestQuery(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	rows, err := DB.Query(context.Background(), "SELECT * FROM xsql")
	a.Nil(err)
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
	err := DB.Insert(context.Background(), &test).Error
	a.Nil(err)
}

func TestEmbeddingInsert(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := TestEmbedding{
		Test1: Test1{
			Id: 0,
		},
		Test2: Test2{
			Foo: "test",
			Bar: time.Now(),
		},
	}
	err := DB.Insert(context.Background(), &test).Error
	a.Nil(err)
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
	err := DB.BatchInsert(context.Background(), &tests).Error
	a.Nil(err)
}

func TestEmbeddingBatchInsert(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	tests := []TestEmbedding{
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
	err := DB.BatchInsert(context.Background(), &tests).Error
	a.Nil(err)
}

func TestEmbeddingUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newDB()
	test := TestEmbedding{
		Test1: Test1{
			Id: 999,
		},
		Test2: Test2{
			Foo: "test_update",
			Bar: time.Now(),
		},
	}
	err := DB.Update(context.Background(), &test, "id = ?", 8).Error
	a.Nil(err)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  8,
		Foo: "test_update_1",
		Bar: time.Now(),
	}
	err := DB.Update(context.Background(), &test, "id = ?", 999).Error
	a.Nil(err)
}

func TestUpdateColumns(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	data := map[string]interface{}{
		"foo": "test_update_2",
	}
	err := DB.Model(&Test{}).Update(context.Background(), data, "id = ?", 8).Error
	a.Nil(err)

	data = map[string]interface{}{
		"foo": timestamppb.Now(),
	}
	err = DB.Model(&Test{}).Update(context.Background(), data, "id = ?", 8).Error

	a.Nil(err)
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
	a.Nil(err)

	err = DB.Model(&test).Update(context.Background(), data, "id = ?", 8).Error
	a.Nil(err)
}

func TestEmbeddingUpdateTagValuesMap(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := TestEmbedding{}
	data, err := xsql.TagValuesMap(DB.Options.Tag, &test,
		xsql.TagValues{
			{&test.Foo, "test_update_4"},
		},
	)
	a.Nil(err)

	err = DB.Model(&test).Update(context.Background(), data, "id = ?", 8).Error
	a.Nil(err)
}

func TestDelete(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	test := Test{
		Id:  8,
		Foo: "test",
		Bar: time.Now(),
	}
	err := DB.Model(&test).Delete(context.Background(), "id = ?", test.Id).Error
	a.Nil(err)

	err = DB.Model(&Test{}).Delete(context.Background(), "id = ?", 8).Error
	a.Nil(err)
}

func TestExec(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	err := DB.Exec(context.Background(), "DELETE FROM xsql WHERE id = ?", 7).Error

	a.Nil(err)
}

func TestFirst(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(context.Background(), &test, "SELECT * FROM ${TABLE}").Error
	a.Nil(err)

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

	var test *TestJsonStructPtr
	err := DB.First(context.Background(), &test, "SELECT * FROM ${TABLE}").Error
	a.Nil(err)

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z","Json":{"Foo":"bar"}}`)
	// bool
	a.Equal(test.Bool, true)
	// enum
	a.IsType(Enum(0), test.Enum)
	a.Equal(Enum(1), test.Enum)
}

func TestFirstEmbedding(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test TestEmbedding
	err := DB.First(context.Background(), &test, "SELECT * FROM ${TABLE}").Error
	a.Nil(err)

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
}

func TestFirstPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(context.Background(), &test, "SELECT foo FROM ${TABLE}").Error
	a.Nil(err)

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":0,"Foo":"v","Bar":"0001-01-01T00:00:00Z"}`)
}

func TestFirstTableKey(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var test Test
	err := DB.First(context.Background(), &test, "SELECT * FROM ${TABLE}").Error
	a.Nil(err)

	b, _ := json.Marshal(test)
	a.Equal(string(b), `{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"}`)
}

func TestFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(context.Background(), &tests, "SELECT * FROM ${TABLE} LIMIT 2").Error
	a.Nil(err)

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
}

func TestFindPtr(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []*TestJsonStructPtr
	err := DB.Find(context.Background(), &tests, "SELECT * FROM ${TABLE} LIMIT 1").Error
	a.Nil(err)

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z","Json":{"Foo":"bar"}}]`)
}

func TestEmbeddingFind(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []TestEmbedding
	err := DB.Find(context.Background(), &tests, "SELECT * FROM ${TABLE} LIMIT 2").Error
	a.Nil(err)

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":1,"Foo":"v","Bar":"2022-04-12T23:50:00Z"},{"Id":2,"Foo":"v1","Bar":"2022-04-13T23:50:00Z"}]`)
}

func TestFindPart(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(context.Background(), &tests, "SELECT foo FROM ${TABLE} LIMIT 2").Error
	a.Nil(err)

	b, _ := json.Marshal(tests)
	a.Equal(string(b), `[{"Id":0,"Foo":"v","Bar":"0001-01-01T00:00:00Z"},{"Id":0,"Foo":"v1","Bar":"0001-01-01T00:00:00Z"}]`)
}

func TestFindTableKey(t *testing.T) {
	a := assert.New(t)

	DB := newDB()

	var tests []Test
	err := DB.Find(context.Background(), &tests, "SELECT * FROM ${TABLE} LIMIT 2").Error
	a.Nil(err)

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
	err := tx.Insert(context.Background(), &test).Error
	a.Nil(err)

	err = tx.Commit()
	a.Nil(err)
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
	err := tx.Insert(context.Background(), &test).Error
	a.Nil(err)

	err = tx.Rollback()
	a.Nil(err)
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
	res := DB.Insert(context.Background(), &test)
	err := res.Error
	a.Nil(err)
	insertId := res.LastInsertId

	// First
	var test2 Test3
	err = DB.First(context.Background(), &test2, "SELECT * FROM ${TABLE} WHERE id = ?", insertId).Error
	a.Nil(err)
	// Timestamp
	a.IsType(&timestamppb.Timestamp{}, test2.Bar)
	a.Equal(test2.Bar.Seconds, now.Seconds)
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
	err := DB.Insert(context.Background(), &test1).Error
	a.Nil(err)

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
	err = DB.Insert(context.Background(), &test2).Error
	a.Nil(err)

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
	err = DB.Insert(context.Background(), &test3).Error
	a.Nil(err)

	test4 := TestJsonSlicePtr{
		Test: Test{
			Id:   0,
			Foo:  "",
			Bar:  time.Time{},
			Bool: false,
			Enum: 0,
		},
		Json: []*JsonItem{{Foo: `bar1`}, {Foo: `bar2`}, {Foo: `bar3`}},
	}
	err = DB.Insert(context.Background(), &test4).Error
	a.Nil(err)
}

func TestFirstPbJsonField(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	var test1 TestJsonStruct
	err := DB.First(context.Background(), &test1, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Nil(err)
	a.NotEmpty(test1.Json)

	var test2 TestJsonStructPtr
	err = DB.First(context.Background(), &test2, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Nil(err)
	a.NotEmpty(test2.Json)

	var test3 TestJsonSlice
	err = DB.First(context.Background(), &test3, "SELECT * FROM ${TABLE} WHERE id = 2").Error
	a.Nil(err)
	a.NotEmpty(test3.Json)

	var test4 TestJsonSlicePtr
	err = DB.First(context.Background(), &test4, "SELECT * FROM ${TABLE} WHERE id = 1006").Error
	a.Nil(err)
	a.NotEmpty(test4.Json)
}

func TestFindPbJsonField(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	var test1 []*TestJsonStruct
	err := DB.Find(context.Background(), &test1, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Nil(err)
	a.NotEmpty(test1)

	var test2 []*TestJsonStructPtr
	err = DB.Find(context.Background(), &test2, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Nil(err)
	a.NotEmpty(test2)

	var test3 []*TestJsonSlice
	err = DB.Find(context.Background(), &test3, "SELECT * FROM ${TABLE} WHERE id = 2").Error
	a.Nil(err)
	a.NotEmpty(test3)

	var test4 []*TestJsonSlicePtr
	err = DB.Find(context.Background(), &test4, "SELECT * FROM ${TABLE} WHERE id = 1006").Error
	a.Nil(err)
	a.NotEmpty(test4)
}

func TestFirstPbStruct(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	var row TestPbStruct
	err := DB.First(context.Background(), &row, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Nil(err)
	a.NotEmpty(&row)

	var row1 testdata.Device
	err = DB.First(context.Background(), &row1, "SELECT * FROM ${TABLE} WHERE id = 1").Error
	a.Contains(err.Error(), "doesn't exist")
}

func TestFindPbStruct(t *testing.T) {
	a := assert.New(t)
	DB := newDB()

	var rows []*TestPbStruct
	err := DB.Find(context.Background(), &rows, "SELECT * FROM ${TABLE} WHERE id < 3").Error
	a.Nil(err)
	a.Len(rows, 2)

	var rows1 []*testdata.Device
	err = DB.Find(context.Background(), &rows1, "SELECT * FROM ${TABLE} WHERE id < 3").Error
	a.Contains(err.Error(), "doesn't exist")
}
