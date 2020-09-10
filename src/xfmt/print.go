package xfmt

import (
    "errors"
    "fmt"
    "reflect"
    "strings"
)

const (
    varFlag = "%v"
)

type pointer struct {
    Format string
    Ptr    uintptr
    Addr   reflect.Value
}

func Print(depth int, args ...interface{}) {
    fmt.Print(Sprintf(depth, format(args...), args...))
}

func Println(depth int, args ...interface{}) {
    fmt.Println(Sprintf(depth, format(args...), args...))
}

func Printf(depth int, format string, args ...interface{}) {
    fmt.Print(Sprintf(depth, format, args...))
}

func Sprint(depth int, args ...interface{}) string {
    return Sprintf(depth, format(args...), args...)
}

func Sprintln(depth int, args ...interface{}) string {
    return Sprintf(depth, format(args...)+"\n", args...)
}

func Sprintf(depth int, format string, args ...interface{}) string {
    // 放在第一行可以起到效验的作用
    str := fmt.Sprintf(format, args...)

    values := []interface{}{}
    for _, arg := range args {
        switch reflect.ValueOf(arg).Kind() {
        case reflect.Ptr, reflect.Struct:
            values = append(values, arg)
            break
        }
    }
    flags := flags(format)
    if len(values) != len(flags) {
        panic(errors.New("Format invalid"))
    }

    pointers := []pointer{}
    for k, val := range values {
        pointers = append(pointers, extract(reflect.ValueOf(val), depth-1, flags[k])...)
    }

    return replace(str, pointers)
}

func format(args ...interface{}) string {
    flags := []string{}
    for i := 0; i < len(args); i++ {
        flags = append(flags, varFlag)
    }
    return strings.Join(flags, " ")
}

func flags(format string) []string {
    fbytes := []byte(format)
    l := len(fbytes) - 1
    flags := []string{}
    for k, v := range fbytes {
        if v == '%' {
            if k+1 <= l {
                switch fbytes[k+1] {
                case 'v':
                    flags = append(flags, "%v")
                    break
                case '+':
                    if k+2 <= l && fbytes[k+2] == 'v' {
                        flags = append(flags, "%+v")
                    }
                    break
                case '#':
                    if k+2 <= l && fbytes[k+2] == 'v' {
                        flags = append(flags, "%#v")
                    }
                    break
                }
            }
        }
    }
    return flags
}

func replace(str string, pointers []pointer) string {
    for _, ptr := range pointers {
        sptr := fmt.Sprintf("0x%x", ptr.Ptr)
        str = strings.Replace(str, sptr, fmt.Sprintf("%s:"+ptr.Format, sptr, ptr.Addr), 1)
    }
    return str
}

func extract(val reflect.Value, level int, format string) []pointer {
    switch val.Kind() {
    case reflect.Ptr:
        val = val.Elem()
        break
    case reflect.Struct:
        break
    default:
        return []pointer{}
    }
    pointers := []pointer{}
    for i := 0; i < val.NumField(); i++ {
        if val.Field(i).Kind() == reflect.Ptr {
            elem := val.Field(i).Elem()
            if !elem.CanAddr() { // 空指针无法寻址
                continue
            }
            if level > 0 {
                pointers = append(pointers, pointer{
                    Format: format,
                    Ptr:    elem.Addr().Pointer(),
                    Addr:   elem.Addr(),
                })
            }
            if level-1 > 0 {
                pointers = append(pointers, extract(elem, level-1, format)...)
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
