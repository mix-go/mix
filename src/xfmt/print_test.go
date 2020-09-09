package xfmt

import (
    "fmt"
    "testing"
)

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

func Test(t *testing.T) {
    l3 := level3{name: "level3"}
    l2 := level2{name: "level2", level3: &l3}
    l1 := level1{name: "level1", level2: &l2, level2_1: &l2}
    fmt.Println(Sprintf(1, "%+v", &l1))
    fmt.Println(Sprintf(1, "%+v", l1))
    fmt.Println(Sprintf(2, "%+v", l1))
    fmt.Println(Sprintf(3, "%+v", l1))
    fmt.Println(Sprintf(100, "%+v", l1))
    fmt.Println(fmt.Sprintf("%+v", l1))
}
