> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XFMT

可以打印结构体嵌套指针地址内部数据的格式化库

Formatting library that can print the internal data of the nested pointer address of the struct

## Overview

在 Golang 中使用 `fmt` 打印结构体时，无法打印指针字段内部的数据结构，导致增加 debug 难度，该库可以解决这个问题。

## Installation

- 安装

```
go get github.com/mix-go/xfmt
```

## Usage

- 支持的方法，与 `fmt` 系统库完全一致

  - `Sprintf(format string, args ...interface{}) string` 
  - `Sprint(args ...interface{}) string` 
  - `Sprintln(args ...interface{}) string` 
  - `Printf(format string, args ...interface{})` 
  - `Print(args ...interface{})` 
  - `Println(args ...interface{})` 

- 支持 `Tag` 忽略某个引用字段

```go
type Foo struct {
    Bar *Bar `xfmt:"-"`
}
```

- 使用

包含指针的结构体

```go
type Level3 struct {
    Name string
}

type Level2 struct {
    Level3 *Level3 `xfmt:"-"`
    Name   string
}

type Level1 struct {
    Name     string
    Level2   *Level2
    Level2_1 *Level2
}
```

创建变量

```go
l3 := Level3{Name: "Level3"}
l2 := Level2{Name: "Level2", Level3: &l3}
l1 := Level1{Name: "Level1", Level2: &l2, Level2_1: &l2}
```

打印对比

- `fmt` 打印

```go
fmt.Println(fmt.Sprintf("%+v", l1))
```

```
{Name:Level1 Level2:0xc00009c500 Level2_1:0xc00009c500}
```

- `xfmt` 打印：其中 Level3 被定义的 tag 忽略，Level2_1 由于和 Level2 是同一个指针因此后面的忽略处理

```go
fmt.Println(xfmt.Sprintf("%+v", l1))
```

```
{Name:Level1 Level2:0xc00009c500:&{Level3:0xc00007f030 Name:Level2} Level2_1:0xc00009c500}
```

动态停用和启用

```go
xfmt.Disable() // 停用后xfmt等同于fmt
xfmt.Enable()
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
