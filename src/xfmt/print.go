package xfmt

import (
    "fmt"
    "reflect"
    "strings"
)

const (
    defaultFormat = "%v"
)

type pointer struct {
    Ptr  uintptr
    Addr reflect.Value
}

func Print(depth int, a interface{}) {
    fmt.Print(Sprintf(depth, defaultFormat, a))
}

func Println(depth int, a interface{}) {
    fmt.Println(Sprintf(depth, defaultFormat+"\n", a))
}

func Printf(depth int, format string, a interface{}) {
    fmt.Print(Sprintf(depth, format, a))
}

func Sprint(depth int, a interface{}) string {
    return Sprintf(depth, defaultFormat, a)
}

func Sprintln(depth int, a interface{}) string {
    return Sprintf(depth, defaultFormat+"\n", a)
}

func Sprintf(depth int, format string, a interface{}) string {
    str := fmt.Sprintf(format, a) // 放在第一行可以起到效验的作用
    pointers := extract(reflect.ValueOf(a), depth-1)
    return replace(str, format, pointers)
}

func replace(str string, format string, pointers []pointer) string {
    for _, ptr := range pointers {
        sptr := fmt.Sprintf("0x%x", ptr.Ptr)
        str = strings.Replace(str, sptr, fmt.Sprintf("%s:"+format, sptr, ptr.Addr), 1)
    }
    return str
}

func extract(val reflect.Value, level int) []pointer {
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    pointers := []pointer{}
    for i := 0; i < val.NumField(); i++ {
        if val.Field(i).Kind() == reflect.Ptr {
            elem := val.Field(i).Elem()
            if level > 0 {
                pointers = append(pointers, pointer{
                    Ptr:  elem.Addr().Pointer(),
                    Addr: elem.Addr(),
                })
            }
            if level-1 > 0 {
                pointers = append(pointers, extract(elem, level-1)...)
            }
        }
    }
    unique := []pointer{}
    for _, p := range pointers {
        find := false
        for _, u := range unique {
            if p.Ptr == u.Ptr {
                find = true
            }
        }
        if !find {
            unique = append(unique, p)
        }
    }
    return unique
}
