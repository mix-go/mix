package xsql

import (
	"database/sql"
	"errors"
	"fmt"
	ora "github.com/sijms/go-ora/v2"
	"reflect"
	"strconv"
	"time"
)

type Fetcher struct {
	R       *sql.Rows
	Log     *Log
	Options *sqlOptions
}

func (t *Fetcher) First(i interface{}) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return errors.New("xsql: argument can only be pointer type")
	}
	rootValue := value.Elem()
	rootType := reflect.TypeOf(i).Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return sql.ErrNoRows
	}
	row := rows[0]
	if err := t.foreach(&row, rootValue, rootType); err != nil {
		return err
	}

	return nil
}

func (t *Fetcher) Find(i interface{}) error {
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return errors.New("xsql: argument can only be pointer type")
	}
	root := value.Elem()
	itemType := root.Type().Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}

	for r := 0; r < len(rows); r++ {
		itemValue := reflect.New(itemType)
		if itemValue.Kind() == reflect.Ptr {
			itemValue = itemValue.Elem()
		}
		if err := t.foreach(&rows[r], itemValue, itemValue.Type()); err != nil {
			return err
		}
		root.Set(reflect.Append(root, itemValue))
	}

	return nil
}

func (t *Fetcher) Rows() ([]Row, error) {
	var debugFunc DebugFunc
	if t.Options.DebugFunc != nil {
		debugFunc = t.Options.DebugFunc
	}

	// 获取列名
	columns, err := t.R.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	var rows []Row

	for t.R.Next() {
		err = t.R.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, value := range values {
			// Here we can check if the value is nil (NULL value)
			if value != nil {
				rowMap[columns[i]] = value
			}
		}

		rows = append(rows, Row{
			v:       rowMap,
			options: t.Options,
		})
	}

	if debugFunc != nil {
		t.Log.RowsAffected = int64(len(rows))
		debugFunc(t.Log)
	}

	return rows, nil
}

type Row struct {
	v       map[string]interface{}
	options *sqlOptions
}

func (t Row) Exist(field string) bool {
	_, ok := t.v[field]
	return ok
}

func (t Row) Get(field string) *RowResult {
	if v, ok := t.v[field]; ok {
		return &RowResult{
			v:       v,
			options: t.options,
		}
	}
	return &RowResult{
		v:       "",
		options: t.options,
	}
}

func (t Row) Value() map[string]interface{} {
	return t.v
}

type RowResult struct {
	v       interface{}
	options *sqlOptions
}

func (t *RowResult) Empty() bool {
	if b, ok := t.v.([]uint8); ok {
		return len(b) == 0
	}
	if s, ok := t.v.(string); ok {
		return len(s) == 0
	}
	if t.v == nil {
		return true
	}
	return false
}

func (t *RowResult) String() string {
	switch reflect.ValueOf(t.v).Kind() {
	case reflect.Int:
		i := t.v.(int)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int8:
		i := t.v.(int8)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int16:
		i := t.v.(int16)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int32:
		i := t.v.(int32)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Int64:
		i := t.v.(int64)
		return strconv.FormatInt(i, 10)
	case reflect.Uint:
		i := t.v.(uint)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint8:
		i := t.v.(uint8)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint16:
		i := t.v.(uint16)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint32:
		i := t.v.(uint32)
		return strconv.FormatInt(int64(i), 10)
	case reflect.Uint64:
		i := t.v.(uint64)
		return strconv.FormatInt(int64(i), 10)
	case reflect.String:
		return t.v.(string)
	default:
		switch v := t.v.(type) {
		case []uint8:
			return string(v)
		case time.Time:
			return v.Format(t.options.TimeLayout)
		}
	}
	return ""
}

func (t *RowResult) Int() int64 {
	switch reflect.ValueOf(t.v).Kind() {
	case reflect.Int:
		i := t.v.(int)
		return int64(i)
	case reflect.Int8:
		i := t.v.(int8)
		return int64(i)
	case reflect.Int16:
		i := t.v.(int16)
		return int64(i)
	case reflect.Int32:
		i := t.v.(int32)
		return int64(i)
	case reflect.Int64:
		i := t.v.(int64)
		return i
	case reflect.Uint:
		i := t.v.(uint)
		return int64(i)
	case reflect.Uint8:
		i := t.v.(uint8)
		return int64(i)
	case reflect.Uint16:
		i := t.v.(uint16)
		return int64(i)
	case reflect.Uint32:
		i := t.v.(uint32)
		return int64(i)
	case reflect.Uint64:
		i := t.v.(uint64)
		return int64(i)
	case reflect.String:
		s := t.v.(string)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		if b, ok := t.v.([]uint8); ok {
			s := string(b)
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return 0
			}
			return i
		}
	}
	return 0
}

func (t *RowResult) Time() time.Time {
	typ := t.Type()
	if typ == "string" || typ == "[]uint8" {
		tt, _ := time.ParseInLocation(t.options.TimeLayout, t.String(), t.options.TimeLocation)
		return tt
	}
	if typ == "time.Time" {
		return t.v.(time.Time)
	}
	if typ == "ora.TimeStamp" {
		return time.Time(t.v.(ora.TimeStamp))
	}
	return time.Time{}
}

func (t *RowResult) Value() interface{} {
	return t.v
}

func (t *RowResult) Type() string {
	return reflect.TypeOf(t.v).String()
}

func (t *Fetcher) foreach(row *Row, value reflect.Value, typ reflect.Type) error {
	for n := 0; n < typ.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			if err := t.foreach(row, fieldValue, fieldValue.Type()); err != nil {
				return err
			}
			continue
		}
		if !fieldValue.CanSet() {
			continue
		}
		tag := value.Type().Field(n).Tag.Get(t.Options.Tag)
		if tag == "-" || tag == "_" {
			continue
		}
		if !row.Exist(tag) {
			continue
		}
		if err := t.mapped(fieldValue, row, tag); err != nil {
			return err
		}
	}
	return nil
}

func (t *Fetcher) mapped(field reflect.Value, row *Row, tag string) (err error) {
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
			if t, e := time.ParseInLocation(t.Options.TimeLayout, res.String(), t.Options.TimeLocation); e == nil {
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
