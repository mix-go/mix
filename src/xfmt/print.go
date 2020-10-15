package xfmt

import (
    "fmt"
    "reflect"
    "strings"
)

const (
    varFlag = "%v"
)

type value struct {
    Arg  interface{}
    Flag string
}

type pointer struct {
    Format string
    Ptr    uintptr
    Addr   reflect.Value
}

// Print
func Print(depth int, args ...interface{}) {
    fmt.Print(Sprintf(depth, format(args...), args...))
}

// Println
func Println(depth int, args ...interface{}) {
    fmt.Println(Sprintf(depth, format(args...), args...))
}

// Printf
func Printf(depth int, format string, args ...interface{}) {
    fmt.Print(Sprintf(depth, format, args...))
}

// Sprint
func Sprint(depth int, args ...interface{}) string {
    return Sprintf(depth, format(args...), args...)
}

// Sprintln
func Sprintln(depth int, args ...interface{}) string {
    return Sprintf(depth, format(args...)+"\n", args...)
}

// Sprintf
func Sprintf(depth int, format string, args ...interface{}) string {
    // 放在第一行可以起到效验的作用
    str := fmt.Sprintf(format, args...)

    pointers := []pointer{}
    for _, v := range values(format, args...) {
        pointers = append(pointers, extract(reflect.ValueOf(v.Arg), depth-1, v.Flag)...)
    }

    // 去重
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

    return replace(str, unique)
}

// 获取全部指针对象的反射
func values(format string, args ...interface{}) []value {
    flags := flags(format)
    values := []value{}
    for k, arg := range args {
        val := reflect.ValueOf(arg)
        switch val.Kind() {
        case reflect.Struct:
            values = append(values, value{arg, flags[k]})
            break
        case reflect.Ptr:
            elem := val.Elem()
            if !elem.CanAddr() {
                continue
            }
            switch elem.Kind() {
            case reflect.Struct:
                values = append(values, value{arg, flags[k]})
                break
            case reflect.Map:
                iter := elem.MapRange()
                for iter.Next() {
                    values = append(values, value{iter.Value().Interface(), flags[k]})
                }
                break
            case reflect.Slice, reflect.Array:
                for i := 0; i < elem.Len(); i++ {
                    values = append(values, value{elem.Index(i).Interface(), flags[k]})
                }
                break
            }
            break
        case reflect.Map:
            iter := val.MapRange()
            for iter.Next() {
                values = append(values, value{iter.Value().Interface(), flags[k]})
            }
            break
        case reflect.Slice, reflect.Array:
            for i := 0; i < val.Len(); i++ {
                values = append(values, value{val.Index(i).Interface(), flags[k]})
            }
            break
        }
    }
    return values
}

// 生成格式
func format(args ...interface{}) string {
    flags := []string{}
    for i := 0; i < len(args); i++ {
        flags = append(flags, varFlag)
    }
    return strings.Join(flags, " ")
}

// 获取全部参数的格式
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

// 替换指针为机构体
func replace(str string, pointers []pointer) string {
    for _, ptr := range pointers {
        sptr := fmt.Sprintf("0x%x", ptr.Ptr)
        str = strings.Replace(str, sptr, fmt.Sprintf("%s:"+ptr.Format, sptr, ptr.Addr), 1)
    }
    return str
}

// 提取指针信息
func extract(val reflect.Value, level int, format string) []pointer {
    pointers := []pointer{}
    switch val.Kind() {
    case reflect.Ptr:
        elem := val.Elem()
        if elem.Kind() != reflect.Struct || !elem.CanAddr() {
            return []pointer{}
        }
        pointers = append(pointers, pointer{
            Format: format,
            Ptr:    elem.Addr().Pointer(),
            Addr:   elem.Addr(),
        })
        val = elem
        break
    case reflect.Struct:
        break
    default:
        return []pointer{}
    }
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
    return pointers
}
