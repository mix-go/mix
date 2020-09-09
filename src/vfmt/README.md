## Mix VFMT

可以打印变量内部指针的格式化库

Formatting library that can print pointers inside variable

## Usage

- 安装

```
go get -u github.com/mix-go/vfmt
```

- 支持的方法

  - `Sprintf(depth int, format string, a interface{}) string` 
  - `Sprint(depth int, a interface{}) string` 
  - `Sprintln(depth int, a interface{}) string` 
  - `Printf(depth int, format string, a interface{}) (n int, err error)` 
  - `Print(depth int, a interface{}) (n int, err error)` 
  - `Println(depth int, a interface{}) (n int, err error)` 

- 使用

包含指针的结构体

```
type level1 struct {
    level2   *level2
    name     string
    level2_1 *level2
}

type level2 struct {
    level3 *level3
    name   string
}

type level3 struct {
    name string
}
```

创建变量

```
l3 := level3{name: "level3"}
l2 := level2{name: "level2", level3: &l3}
l1 := level1{name: "level1", level2: &l2, level2_1: &l2}
```

打印

```
fmt.Println(vfmt.Sprintf(1, "%+v", l1))
fmt.Println(vfmt.Sprintf(2, "%+v", l1))
fmt.Println(vfmt.Sprintf(3, "%+v", l1))
```

```
{level2:0xc00000c0e0 name:level1 level2_1:0xc00000c0e0}
{level2:0xc00000c0e0=&{level3:0xc00003e480 name:level2} name:level1 level2_1:0xc00000c0e0}
{level2:0xc00000c0e0=&{level3:0xc00003e480=&{name:level3} name:level2} name:level1 level2_1:0xc00000c0e0}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
