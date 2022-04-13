package xsql

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
)

type Row struct {
	m map[string]interface{}
}

func (t *Row) Exist(field string) bool {
	_, ok := t.m[field]
	return ok
}

func (t *Row) Get(field string) *Result {
	if v, ok := t.m[field]; ok {
		return &Result{v: v}
	}
	return &Result{v: ""}
}

func (t *Row) Value() map[string]interface{} {
	return t.m
}

type Result struct {
	v interface{}
}

func (t *Result) String() string {
	if s, ok := t.v.(string); ok {
		return s
	}
	if b, ok := t.v.([]uint8); ok {
		return string(b)
	}
	return ""
}

func (t *Result) Int() int64 {
	s := t.String()
	if s == "" {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func (t *Result) Value() interface{} {
	return t.v
}

func (t *Result) Type() string {
	return reflect.TypeOf(t.v).String()
}

type Fetcher struct {
	r *sql.Rows
}

func Fetch(r *sql.Rows) *Fetcher {
	return &Fetcher{r: r}
}

func (t *Fetcher) First(i interface{}) (err error) {
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Ptr:
	default:
		return errors.New("argument can only be pointer type")
	}
	root := value.Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("rows is empty")
	}
	row := rows[0]

	for n := 0; n < root.NumField(); n++ {
		field := root.Field(n)
		if !field.CanSet() {
			continue
		}
		tag := root.Type().Field(n).Tag.Get("xsql")
		if tag == "-" || tag == "_" {
			continue
		}
		if !row.Exist(tagName(tag)) {
			continue
		}
		if err = mapped(field, row, tag); err != nil {
			return
		}
	}

	return
}

func (t *Fetcher) Find(i interface{}) error {
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Ptr:
	default:
		return errors.New("argument can only be pointer type")
	}
	root := value.Elem()

	rows, err := t.Rows()
	if err != nil {
		return err
	}

	for r := 0; r < root.Len(); r++ {
		obj := root.Index(r).Interface()
		val := reflect.ValueOf(obj)
		for n := 0; n < val.NumField(); n++ {
			field := val.Field(n)
			if !field.CanSet() {
				continue
			}
			tag := val.Type().Field(n).Tag.Get("xsql")
			if !rows[r].Exist(tag) {
				continue
			}
			field.Set(reflect.ValueOf(rows[r].Get(tag).Value()))
		}
	}

	return nil
}

func (t *Fetcher) Rows() ([]*Row, error) {
	// 获取列名
	columns, err := t.r.Columns()
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
	var rows []*Row

	for t.r.Next() {
		err = t.r.Scan(scanArgs...)
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

		rows = append(rows, &Row{m: rowMap})
	}

	return rows, nil
}
