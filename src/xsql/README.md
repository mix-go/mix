> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XRows

database/sql 查询与映射

database/sql query and mapper

## Installation

```
go get github.com/mix-go/xsql
```

## 查询

不需要绑定结构体，自由获取每个字段的值。

```go
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

rows, err := Query(db, "SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}
bar := rows[0].Get("bar").String()
fmt.Println(bar)
```

## 映射

```go
type Test struct {
	Id  int       `xsql:"id"`
	Foo string    `xsql:"foo"`
	Bar time.Time `xsql:"bar"`
}
```

- `First()` 映射第一行

```go
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

rows, err := db.Query("SELECT * FROM xsql limit 1")
if err != nil {
    log.Fatal(err)
}

var test Test
err = First(rows, &test)
if err != nil {
    log.Fatal(err)
}
```

- `Find()` 映射全部行

```go
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
if err != nil {
    log.Fatal(err)
}

rows, err := db.Query("SELECT * FROM xsql")
if err != nil {
    log.Fatal(err)
}
var tests []Test
err = Find(rows, &tests)
if err != nil {
    log.Fatal(err)
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
