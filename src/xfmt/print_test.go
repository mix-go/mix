package xfmt

import (
    "fmt"
    "testing"
)

type level1 struct {
    level2   *level2
    name     string
    level2_1 *level2
    level4   *level4
    level4_1 level4
}

type level2 struct {
    level3 *level3
    name   string
}

type level3 struct {
    name string
}

type level4 struct {
    level5 *level5
    name   string
}

type level5 struct {
    name string
}

func TestRun(t *testing.T) {
    l5 := level5{name: "level5"}
    l4 := level4{name: "level4", level5: &l5}
    l3 := level3{name: "level3"}
    l2 := level2{name: "level2", level3: &l3}
    l1 := level1{name: "level1", level2: &l2, level2_1: &l2, level4: &l4, level4_1: l4}

    fmt.Println(fmt.Sprintf("%+v", l1))

    fmt.Println(Sprintf(1, "%+v", &l1))
    fmt.Println(Sprintf(1, "%+v", l1))
    fmt.Println(Sprintf(2, "%+v", l1))
    fmt.Println(Sprintf(3, "%+v", l1))
    fmt.Println(Sprintf(100, "%+v", l1))
    fmt.Println(Sprintf(100, "%v And %+v And %#v And %s", l3, l4, &l5, "str1"))
}

func TestMultiple(t *testing.T) {
    l5 := level5{name: "level5"}
    l4 := level4{name: "level4", level5: &l5}

    fmt.Println(l4, &l5)

    Print(2, l4, &l5)
    println("")
    Println(2, l4, &l5)
    Printf(2, "%v %v\n", l4, &l5)
    println(Sprint(2, l4, &l5))
    print(Sprintln(2, l4, &l5))
    print(Sprintf(2, "%v %v\n", l4, &l5))
}

func TestMap(t *testing.T) {
    m := map[string]*level5{}
    m["foo"] = &level5{}
    print(Sprintf(2, "%v\n", m))
    print(Sprintf(2, "%v\n", &m))
}

func TestArray(t *testing.T) {
    a := []*level5{}
    a = append(a, &level5{})
    print(Sprintf(2, "%v\n", a))
    print(Sprintf(2, "%v\n", &a))
}
