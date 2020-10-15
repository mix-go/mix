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

// 生成格式
func format(args ...interface{}) string {
    flags := []string{}
    for i := 0; i < len(args); i++ {
        flags = append(flags, varFlag)
    }
    return strings.Join(flags, " ")
}

// 获取全部指针对象的反射
func values(format string, args ...interface{}) []value {
    flags := flags(format)
    result := []value{}
    for k, arg := range args {
        val := reflect.ValueOf(arg)
        switch val.Kind() {
        case reflect.Struct:
            result = append(result, value{arg, flags[k]})
            break
        case reflect.Ptr:
            if !val.Elem().CanAddr() {
                continue
            }
            result = append(result, value{arg, flags[k]})
            break
        case reflect.Map:
            iter := val.MapRange()
            for iter.Next() {
                result = append(result, value{iter.Value().Interface(), flags[k]})
            }
            break
        case reflect.Slice, reflect.Array:
            for i := 0; i < val.Len(); i++ {
                result = append(result, value{val.Index(i).Interface(), flags[k]})
            }
            break
        }
    }
    return result
}

// 获取全部参数的格式
func flags(format string) []string {
    fbytes := []byte(format)
    l := len(fbytes) - 1
    result := []string{}
    for k, v := range fbytes {
        if v == '%' {
            if k+1 <= l {
                switch fbytes[k+1] {
                case 'v':
                    result = append(result, "%v")
                    break
                case '+':
                    if k+2 <= l && fbytes[k+2] == 'v' {
                        result = append(result, "%+v")
                    }
                    break
                case '#':
                    if k+2 <= l && fbytes[k+2] == 'v' {
                        result = append(result, "%#v")
                    }
                    break
                }
            }
        }
    }
    return result
}

// 提取指针信息
func extract(val reflect.Value, depth int, format string) []pointer {
    pointers := []pointer{}
    if depth <= 0 {
        return pointers
    }
    switch val.Kind() {
    case reflect.Ptr:
        elem := val.Elem()
        if !elem.CanAddr() {
            return pointers
        }
        pointers = append(pointers, pointer{
            Format: format,
            Ptr:    elem.Addr().Pointer(),
            Addr:   elem.Addr(),
        })
        for _, v := range values(format, elem.Interface()) {
            pointers = append(pointers, extract(reflect.ValueOf(v.Arg), depth-1, v.Flag)...)
        }
        break
    case reflect.Struct:
        for i := 0; i < val.NumField(); i++ {
            if !val.Field(i).CanInterface() {
                continue
            }
            for _, v := range values(format, val.Field(i).Interface()) {
                pointers = append(pointers, extract(reflect.ValueOf(v.Arg), depth-1, v.Flag)...)
            }
        }
        break
    case reflect.Map:
        for _, v := range values(format, val.Interface()) {
            pointers = append(pointers, extract(reflect.ValueOf(v.Arg), depth-1, v.Flag)...)
        }
    case reflect.Slice, reflect.Array:
        for _, v := range values(format, val.Interface()) {
            pointers = append(pointers, extract(reflect.ValueOf(v.Arg), depth-1, v.Flag)...)
        }
    }
    return pointers
}

// 替换指针为机构体
func replace(str string, pointers []pointer) string {
    for _, ptr := range pointers {
        ptrString := fmt.Sprintf("0x%x", ptr.Ptr)
        str = strings.Replace(str, ptrString, fmt.Sprintf("%s:"+ptr.Format, ptrString, ptr.Addr), 1)
    }
    return str
}
