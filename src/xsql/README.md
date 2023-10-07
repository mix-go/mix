> OpenMix Production: [https://openmix.org](https://openmix.org/mix-go)

## Mix XSQL

A lightweight database based on database/sql, feature complete and supports any database driver.

## Installation

```
go get github.com/mix-go/xsql
```

## Initialization

- MySQL initialization, using [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) driver.

```go
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

DB := xsql.New(db)
```

- Oracle initialization, using [sijms/go-ora/v2](https://github.com/sijms/go-ora) driver (no need to install instantclient).

```go
import _ "github.com/sijms/go-ora/v2"

db, err := sql.Open("oracle", "oracle://root:123456@127.0.0.1:1521/orcl")
if err != nil {
    log.Fatal(err)
}

DB := xsql.New(db, xsql.UseOracle())
```

- [xorm#drivers](https://github.com/go-xorm/xorm#drivers-support) These drivers are also supported

## Query

You can use it like a scripting language, not binding the struct, directly and freely get the value of each field.

> Oracle field, table name needs to be uppercase

```go
rows, err := DB.Query("SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}

id  := rows[0].Get("id").Int()
foo := rows[0].Get("foo").String()
bar := rows[0].Get("bar").Time() // time.Time
val := rows[0].Get("bar").Value() // interface{}
```

### Mapping

Of course, you can also map usage like `gorm`, `xorm`.

> Oracle field, table name needs to be uppercase

```go
type Test struct {
    Id  int       `xsql:"id"`
    Foo string    `xsql:"foo"`
    Bar time.Time `xsql:"bar"` // oracle uses goora.TimeStamp
}

func (t Test) TableName() string {
    return "tableName"
}
```

### First()

Map the first row

> Oracle placeholder needs to be modified to :id

```go
var test Test
err := DB.First(&test, "SELECT * FROM ${TABLE} WHERE id = ?", 1)
if err != nil {
    log.Fatal(err)
}
```

### Find()

Map all rows

```go
var tests []Test
err := DB.Find(&tests, "SELECT * FROM ${TABLE}")
if err != nil {
    log.Fatal(err)
}
```

## Insert

### Insert()

```go
test := Test{
    Id:  0,
    Foo: "test",
    Bar: time.Now(),
}
res, err := DB.Insert(&test)
if err != nil {
    log.Fatal(err)
}
```

### BatchInsert()

```go
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
res, err := DB.BatchInsert(&tests)
if err != nil {
    log.Fatal(err)
}
```

## Update

> Oracle placeholder needs to be modified to :id

```go
test := Test{
    Id:  10,
    Foo: "update",
    Bar: time.Now(),
}
res, err := DB.Update(&test, "id = ?", 10)
```

## Delete

> Oracle placeholder needs to be modified to :id

```go
test := Test{
    Id:  10,
    Foo: "test",
    Bar: time.Now(),
}
res, err := DB.Delete(&test, "DELETE FROM ${TABLE} WHERE id = ?", test.Id)
```

## Exec

Use `Exec()` to manually execute the delete, you can also manually execute the update operation.

> Oracle placeholder needs to be modified to :id

```go
res, err := DB.Exec("DELETE FROM xsql WHERE id = ?", 10)
```

## Transaction

```go
tx, err := DB.Begin()
if err != nil {
    log.Fatal(err)
}
test := Test{
    Id:  0,
    Foo: "test",
    Bar: time.Now(),
}
res, err := tx.Insert(&test)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}
tx.Commit()
```

## Configuration

You can pass the following configuration object in the `xsql.New()` method

- Default to mysql mode, when switching to oracle, you need to [modify the configuration](https://github.com/mix-go/mix/blob/master/src/xsql/dbora_test.go#L24)
- `Insert()`, `BatchInsert()` can pass in configuration during execution to override insert related configuration, such as modifying InsertKey to REPLACE INTO

```go
type sqlOptions struct {
    // Default: xsql
    Tag string
    
    // Default: INSERT INTO
    InsertKey string

    // Default: ${TABLE}
    TableKey string
    
    // Default: ?
    // oracle can be configured as :%d
    Placeholder string
    
    // Default: `
    // oracle can be configured as "
    ColumnQuotes string
    
    // Default: == DefaultTimeLayout
    TimeLayout string

    // Default: == DefaultTimeFunc
    // oracle can modify this closure to add TO_TIMESTAMP
    TimeFunc TimeFunc

    // Global debug SQL
    DebugFunc DebugFunc
}
```

## Log

Pass in the configuration `DebugFunc` when using the `xsql.New()` method, you can print SQL information using any log library here.

```go
DB := xsql.New(
    db,
    xsql.WithDebugFunc(func(l *xsql.Log) {
        log.Println(l)
    }),
)
```

The log object contains the following fields

```go
type Log struct {
	Duration     time.Duration `json:"duration"`
	SQL          string        `json:"sql"`
	Bindings     []interface{} `json:"bindings"`
	RowsAffected int64         `json:"rowsAffected"`
	Error        error         `json:"error"`
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
