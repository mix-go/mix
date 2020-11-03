> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XFMT

可以打印内部指针的格式化库

A formatting library that prints internal Pointers

## Overview

在 go 中使用 `fmt` 打印结构体时，无法打印指针字段内部的数据结构，导致增加 debug 难度，该库可以解决这个问题，并支持设定打印的深度。

## Installation

- 安装

```
go get -u github.com/mix-go/xfmt
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

```
type Foo struct {
    Bar *Bar `xfmt:"-"`
}
```

- 使用

包含指针的结构体

```
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

```
l3 := Level3{Name: "Level3"}
l2 := Level2{Name: "Level2", Level3: &l3}
l1 := Level1{Name: "Level1", Level2: &l2, Level2_1: &l2}
```

打印

```
fmt.Println(xfmt.Sprintf("%+v", l1))
```

```
{Name:Level1 Level2:0xc00009c500:&{Level3:0xc00007f030 Name:Level2} Level2_1:0xc00009c500}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
