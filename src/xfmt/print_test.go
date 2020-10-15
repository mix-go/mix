package xfmt

import (
    "fmt"
    "testing"
)

type level1 struct {
    Level2   *level2
    Name     string
    Level2_1 *level2
    Level4   *level4
    Level4_1 level4
}

type level2 struct {
    Level3 *level3
    Name   string
}

type level3 struct {
    Name string
}

type level4 struct {
    Level5 *level5
    Name   string
}

type level5 struct {
    Name string
}

type level6 struct {
    N   string
    Map map[string]*level5
    Ary []*level5
}

func TestRun(t *testing.T) {
    l5 := level5{Name: "Level5"}
    l4 := level4{Name: "Level4", Level5: &l5}
    l3 := level3{Name: "Level3"}
    l2 := level2{Name: "Level2", Level3: &l3}
    l1 := level1{Name: "level1", Level2: &l2, Level2_1: &l2, Level4: &l4, Level4_1: l4}

    fmt.Println(fmt.Sprintf("%+v", l1))

    fmt.Println(Sprintf(1, "%+v", &l1))
    fmt.Println(Sprintf(1, "%+v", l1))
    fmt.Println(Sprintf(2, "%+v", l1))
    fmt.Println(Sprintf(3, "%+v", l1))
    fmt.Println(Sprintf(100, "%+v", l1))
    fmt.Println(Sprintf(100, "%v And %+v And %#v And %s", l3, l4, &l5, "str1"))
}

func TestMultiple(t *testing.T) {
    l5 := level5{Name: "Level5"}
    l4 := level4{Name: "Level4", Level5: &l5}

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
    m["bar"] = &level5{}
    print(Sprintf(2, "%v\n", m))
    print(Sprintf(2, "%v\n", &m))

    m2 := map[string]level5{}
    m2["foo"] = level5{}
    m2["bar"] = level5{}
    print(Sprintf(2, "%v\n", m2))
    print(Sprintf(2, "%v\n", &m2))
}

func TestArray(t *testing.T) {
    a := []*level5{}
    a = append(a, &level5{})
    a = append(a, &level5{})
    print(Sprintf(2, "%v\n", a))
    print(Sprintf(2, "%v\n", &a))

    a2 := []level5{}
    a2 = append(a2, level5{})
    a2 = append(a2, level5{})
    print(Sprintf(2, "%v\n", a2))
    print(Sprintf(2, "%v\n", &a2))
}

func TestStructMapArray(t *testing.T) {
    m := map[string]*level5{}
    m["foo"] = &level5{Name: "Level5"}
    m["bar"] = &level5{Name: "Level5"}

    a := []*level5{}
    a = append(a, &level5{Name: "Level5"})
    a = append(a, &level5{Name: "Level5"})

    x := level6{
        N:   "level6",
        Map: m,
        Ary: a,
    }
    print(Sprintf(3, "%v\n", x))
    print(Sprintf(3, "%v\n", &x))
}
