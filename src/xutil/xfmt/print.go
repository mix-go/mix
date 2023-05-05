package xfmt

import (
	"fmt"
	"reflect"
	"strings"
)

const varFormat = "%v"

var status = true

func Enable() {
	status = true
}

func Disable() {
	status = false
}

type value struct {
	Arg  interface{}
	Flag string
}

type pointer struct {
	Format string
	Ptr    uintptr
	Addr   reflect.Value
}

// Print 打印
func Print(args ...interface{}) {
	fmt.Print(Sprintf(format(args...), args...))
}

// Println 打印并换行
func Println(args ...interface{}) {
	fmt.Println(Sprintf(format(args...), args...))
}

// Printf 打印并格式化
func Printf(format string, args ...interface{}) {
	fmt.Print(Sprintf(format, args...))
}

// Sprint 打印并返回
func Sprint(args ...interface{}) string {
	return Sprintf(format(args...), args...)
}

// Sprintln 打印换行并返回
func Sprintln(args ...interface{}) string {
	return Sprintf(format(args...)+"\n", args...)
}

// Sprintf 格式化并返回
func Sprintf(format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...) // 放在第一行可以起到效验的作用

	// 停用
	if !status {
		return str
	}

	var pointers []pointer
	for _, v := range values(format, args...) {
		pointers = append(pointers, extract(reflect.ValueOf(v.Arg), v.Flag)...)
	}

	// 去重
	var unique []pointer
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
	var formats []string
	for i := 0; i < len(args); i++ {
		formats = append(formats, varFormat)
	}
	return strings.Join(formats, " ")
}

// 获取全部指针对象的反射
func values(format string, args ...interface{}) []value {
	var result []value
	for _, v := range filter(format, args...) {
		val := reflect.ValueOf(v.Arg)
		switch val.Kind() {
		case reflect.Struct:
			result = append(result, v)
			break
		case reflect.Ptr:
			if !val.Elem().CanAddr() {
				continue
			}
			result = append(result, v)
			break
		case reflect.Map:
			iter := val.MapRange()
			for iter.Next() {
				result = append(result, value{iter.Value().Interface(), v.Flag})
			}
			break
		case reflect.Slice, reflect.Array:
			for i := 0; i < val.Len(); i++ {
				result = append(result, value{val.Index(i).Interface(), v.Flag})
			}
			break
		}
	}
	return result
}

// 过滤无需解析的参数
func filter(format string, args ...interface{}) []value {
	fb := []byte(format)
	next := len(fb) - 1
	loc := -1
	var result []value
	for k, v := range fb {
		if v == '%' {
			loc += 1
			if k+1 <= next {
				switch fb[k+1] {
				case 'v':
					result = append(result, value{
						Arg:  args[loc],
						Flag: "%v",
					})
					break
				case '+':
					if k+2 <= next && fb[k+2] == 'v' {
						result = append(result, value{
							Arg:  args[loc],
							Flag: "%+v",
						})
					}
					break
				case '#':
					if k+2 <= next && fb[k+2] == 'v' {
						result = append(result, value{
							Arg:  args[loc],
							Flag: "%#v",
						})
					}
					break
				}
			}
		}
	}
	return result
}

// 提取指针信息
func extract(val reflect.Value, format string) []pointer {
	var pointers []pointer
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
			pointers = append(pointers, extract(reflect.ValueOf(v.Arg), v.Flag)...)
		}
		break
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if !val.Field(i).CanInterface() {
				continue
			}
			tag := val.Type().Field(i).Tag.Get("xfmt")
			if tag == "-" || tag == "_" {
				continue
			}
			for _, v := range values(format, val.Field(i).Interface()) {
				pointers = append(pointers, extract(reflect.ValueOf(v.Arg), v.Flag)...)
			}
		}
		break
	case reflect.Map, reflect.Slice, reflect.Array:
		for _, v := range values(format, val.Interface()) {
			pointers = append(pointers, extract(reflect.ValueOf(v.Arg), v.Flag)...)
		}
		break
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
