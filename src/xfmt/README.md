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

- 支持的方法

  - `Sprintf(depth int, format string, args ...interface{}) string` 
  - `Sprint(depth int, args ...interface{}) string` 
  - `Sprintln(depth int, args ...interface{}) string` 
  - `Printf(depth int, format string, args ...interface{})` 
  - `Print(depth int, args ...interface{})` 
  - `Println(depth int, args ...interface{})` 

- 使用

包含指针的结构体

```
type Level1 struct {
    Level2   *Level2
    Name     string
    Level2_1 *Level2
}

type Level2 struct {
    Level3 *Level3
    Name   string
}

type Level3 struct {
    Name string
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
fmt.Println(xfmt.Sprintf(1, "%+v", l1))
fmt.Println(xfmt.Sprintf(2, "%+v", l1))
fmt.Println(xfmt.Sprintf(3, "%+v", l1))
```

```
{Level2:0xc00000c0e0 Name:Level1 Level2_1:0xc00000c0e0}
{Level2:0xc00000c0e0:&{Level3:0xc0000404d0 Name:Level2} Name:Level1 Level2_1:0xc00000c0e0}
{Level2:0xc00000c0e0:&{Level3:0xc0000404d0:&{Name:Level3} Name:Level2} Name:Level1 Level2_1:0xc00000c0e0}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
