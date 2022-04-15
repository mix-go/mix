package xsql

import (
	"fmt"
	"reflect"
	"time"
)

func mapped(field reflect.Value, row Row, tag string) (err error) {
	res := row.Get(tag)
	v := res.Value()

	switch field.Kind() {
	case reflect.Int:
		v = int(res.Int())
		break
	case reflect.Int8:
		v = int8(res.Int())
		break
	case reflect.Int16:
		v = int16(res.Int())
		break
	case reflect.Int32:
		v = int32(res.Int())
		break
	case reflect.Int64:
		v = res.Int()
		break
	case reflect.Uint:
		v = uint(res.Int())
		break
	case reflect.Uint8:
		v = uint8(res.Int())
		break
	case reflect.Uint16:
		v = uint16(res.Int())
		break
	case reflect.Uint32:
		v = uint32(res.Int())
		break
	case reflect.Uint64:
		v = uint64(res.Int())
		break
	case reflect.String:
		v = res.String()
		break
	default:
		if !res.Empty() &&
			field.Type().String() == "time.Time" &&
			reflect.ValueOf(v).Type().String() != "time.Time" {
			if t, e := time.ParseInLocation(TimeParselayout, res.String(), time.Local); e == nil {
				v = t
			} else {
				return fmt.Errorf("time parse fail for field %s: %v", tag, e)
			}
		}
	}

	// 追加异常信息
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("type mismatch for field %s: %v", tag, e)
		}
	}()
	field.Set(reflect.ValueOf(v))

	return
}
