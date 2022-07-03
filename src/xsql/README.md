> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XSQL

database/sql 标准库的查询与映射，支持任何数据库驱动。

## Installation

```
go get github.com/mix-go/xsql
```

## 初始化

```go
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

DB := xsql.New(db)
```

## 查询

不需要绑定结构体，自由获取每个字段的值。

```go
rows, log, err := DB.Query(db, "SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}

bar := rows[0].Get("bar").String()
fmt.Println(bar)
```

### 映射

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
log, err = DB.First(&test, "SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}
```

### `Find()`

映射全部行

```go
var tests []Test
log, err = DB.Find(&tests, "SELECT * FROM xsql")
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
res, log, err := DB.Insert(&test)
if err != nil {
    log.Fatal(err)
}
```

### `BatchInsert()`

```go
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
res, log, err := DB.BatchInsert(&test)
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
res, log, err := DB.Update(&test, "id = ?", 10)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
