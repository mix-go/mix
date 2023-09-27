package xsql_test

import (
	"database/sql"
	"fmt"
	"github.com/mix-go/xsql"
	ora "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func newOracleDB() *xsql.DB {
	db, err := sql.Open("oracle", "oracle://root:123456@127.0.0.1:1521/orcl")
	if err != nil {
		log.Fatal(err)
	}
	return xsql.New(
		db,
		xsql.WithDebugFunc(func(l *xsql.Log) {
			log.Println(l)
		}),
		xsql.SwitchToOracle(),
	)
}

func TestOracleClear(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	_, err := DB.Exec("DELETE FROM XSQL WHERE ID > 2")

	a.Empty(err)
}

func TestOracleQuery(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	rows, err := DB.Query("SELECT * FROM XSQL WHERE ROWNUM <= 2")
	if err != nil {
		log.Fatal(err)
	}
	bar := rows[0].Get("BAR").Time().Format(xsql.DefaultOptions.TimeLayout)
	fmt.Println(bar)

	a.Equal(bar, "2022-04-14 23:49:48")
}

type TestOracle struct {
	Id  int           `xsql:"ID"`
	Foo string        `xsql:"FOO"`
	Bar ora.TimeStamp `xsql:"BAR"`
}

func (t TestOracle) TableName() string {
	return "XSQL"
}

func TestOracleInsert(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	test := TestOracle{
		Id:  3,
		Foo: "test",
		Bar: ora.TimeStamp(time.Now()),
	}
	_, err := DB.Insert(&test)

	a.Empty(err)
}

// oracle 不支持批量插入
func _TestOracleBatchInsert(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	tests := []TestOracle{
		{
			Id:  4,
			Foo: "test",
			Bar: ora.TimeStamp(time.Now()),
		},
		{
			Id:  5,
			Foo: "test",
			Bar: ora.TimeStamp(time.Now()),
		},
	}
	_, err := DB.BatchInsert(&tests)

	a.Empty(err)
}

func TestOracleUpdate(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	test := TestOracle{
		Id:  999,
		Foo: "test update",
		Bar: ora.TimeStamp(time.Now()),
	}
	_, err := DB.Update(&test, "id = :id", 3)

	a.Empty(err)
}

func TestOracleExec(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	_, err := DB.Exec("DELETE FROM XSQL WHERE ID = :id", 999)

	a.Empty(err)
}

func TestOracleFirst(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	var test TestOracle
	err := DB.First(&test, "SELECT * FROM XSQL")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", test), "{Id:1 Foo:v Bar:{wall:0 ext:63785576988 loc:<nil>}}")
}

func TestOracleFirstPart(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	var test TestOracle
	err := DB.First(&test, "SELECT foo FROM XSQL")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", test), "{Id:0 Foo:v Bar:{wall:0 ext:0 loc:<nil>}}")
}

func TestOracleFind(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	var tests []TestOracle
	err := DB.Find(&tests, "SELECT * FROM XSQL WHERE ROWNUM <= 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:1 Foo:v Bar:{wall:0 ext:63785576988 loc:<nil>}} {Id:2 Foo:v1 Bar:{wall:0 ext:63785577000 loc:<nil>}}]`)
}

func TestOracleFindPart(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	var tests []TestOracle
	err := DB.Find(&tests, "SELECT foo FROM XSQL WHERE ROWNUM <= 2")
	if err != nil {
		log.Fatal(err)
	}

	a.Equal(fmt.Sprintf("%+v", tests), `[{Id:0 Foo:v Bar:{wall:0 ext:0 loc:<nil>}} {Id:0 Foo:v1 Bar:{wall:0 ext:0 loc:<nil>}}]`)
}

func TestOracleTxCommit(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	tx, _ := DB.Begin()

	test := TestOracle{
		Id:  999,
		Foo: "test",
		Bar: ora.TimeStamp(time.Now()),
	}
	_, err := tx.Insert(&test)
	a.Empty(err)

	err = tx.Commit()
	a.Empty(err)
}

func TestOracleTxRollback(t *testing.T) {
	a := assert.New(t)

	DB := newOracleDB()

	tx, _ := DB.Begin()

	test := TestOracle{
		Id:  998,
		Foo: "test",
		Bar: ora.TimeStamp(time.Now()),
	}
	_, err := tx.Insert(&test)
	a.Empty(err)

	err = tx.Rollback()
	a.Empty(err)
}
