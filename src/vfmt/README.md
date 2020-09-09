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
  - `Printf(depth int, format string, a interface{})` 
  - `Print(depth int, a interface{})` 
  - `Println(depth int, a interface{})` 

- 使用

包含指针的结构体

```
type Level1 struct {
    Level2   *Level2
    name     string
    Level2_1 *Level2
}

type Level2 struct {
    Level3 *Level3
    name   string
}

type Level3 struct {
    name string
}
```

创建变量

```
l3 := Level3{name: "Level3"}
l2 := Level2{name: "Level2", Level3: &l3}
l1 := Level1{name: "Level1", Level2: &l2, Level2_1: &l2}
```

打印

```
fmt.Println(vfmt.Sprintf(1, "%+v", l1))
fmt.Println(vfmt.Sprintf(2, "%+v", l1))
fmt.Println(vfmt.Sprintf(3, "%+v", l1))
```

```
{Level2:0xc00000c0e0 name:Level1 Level2_1:0xc00000c0e0}
{Level2:0xc00000c0e0=&{Level3:0xc00003e480 name:Level2} name:Level1 Level2_1:0xc00000c0e0}
{Level2:0xc00000c0e0=&{Level3:0xc00003e480=&{name:Level3} name:Level2} name:Level1 Level2_1:0xc00000c0e0}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
