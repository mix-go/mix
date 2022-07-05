> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XSQL

基于 database/sql 的轻量数据库，功能完备且支持任何数据库驱动。

A lightweight database based on database/sql

## 安装

```
go get github.com/mix-go/xsql
```

推荐使用以下数据库驱动，其他驱动基本也都支持

- Mysql `_ "github.com/go-sql-driver/mysql"`
- Oracle `_ "github.com/sijms/go-ora/v2"` 无需安装 instantclient
- [xorm#drivers](https://github.com/go-xorm/xorm#drivers-support) 这些驱动都支持

## 初始化

```go
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

DB := xsql.New(db)
```

## 查询

可以像脚本语言一样使用，不绑定结构体，直接自由获取每个字段的值。

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

### 映射

当然你也可以像 `gorm`, `xorm` 一样映射使用。

```go
type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}

func (t Test) TableName() string {
    return "tableName"
}
```

### `First()`

映射第一行

```go
var test Test
err := DB.First(&test, "SELECT * FROM xsql WHERE id = ?", 1)
if err != nil {
    log.Fatal(err)
}
```

### `Find()`

映射全部行

```go
var tests []Test
err := DB.Find(&tests, "SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}
```

## 插入

### `Insert()`

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

### `BatchInsert()`

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

## 更新

```go
test := Test{
    Id:  10,
    Foo: "update",
    Bar: time.Now(),
}
res, err := DB.Update(&test, "id = ?", 10)
```

## 删除

采用 `Exec()` 手动执行删除，也可手动执行更新操作。

```go
res, err := DB.Exec("DELETE FROM xsql WHERE id = ?", 10)
```

## 事务

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

## 配置

在 `xsql.New()` 方法中可以传入以下配置对象

- 默认为 mysql 模式，当切换到 oracle 时，需要修改 `Placeholder`、`ColumnQuotes` 配置
- `Insert()`、`BatchInsert()` 可在执行时传入配置，覆盖 insert 相关的配置，比如将 InsertKey 修改为 REPLACE INTO

```go
type Options struct {
    // 默认: INSERT INTO
    InsertKey string
    
    // 默认: ?
    // oracle 可配置为 :%d
    Placeholder string
    
    // 默认：`
    // oracle 可配置为 "
    ColumnQuotes string
    
    // 默认：== DefaultTimeLayout
    TimeLayout string
	
    // 全局 debug SQL
    DebugFunc DebugFunc
}
```

## 日志

在 `xsql.New()` 方法时传入配置 `DebugFunc`，可以在这里使用任何日志库打印SQL信息。

```go
opts := Options{
    DebugFunc: func(l *Log) {
        log.Println(l)
    },
}
DB := New(db, opts)
```

日志对象包含以下字段

```go
type Log struct {
    Time     time.Duration `json:"time"`
    SQL      string        `json:"sql"`
    Bindings []interface{} `json:"bindings"`
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
