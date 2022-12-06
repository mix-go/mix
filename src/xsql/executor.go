package xsql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Table interface {
	TableName() string
}

type executor struct {
	Executor
}

func (t *executor) Insert(data interface{}, opts *Options) (sql.Result, error) {
	insertKey := "INSERT INTO"
	if opts.InsertKey != "" {
		insertKey = opts.InsertKey
	}
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeLayout := DefaultTimeLayout
	if opts.TimeLayout != "" {
		timeLayout = opts.TimeLayout
	}
	timeFunc := DefaultTimeFunc
	if opts.TimeFunc != nil {
		timeFunc = opts.TimeFunc
	}
	columnQuotes := "`"
	if opts.ColumnQuotes != "" {
		columnQuotes = opts.ColumnQuotes
	}
	var debugFunc DebugFunc
	if opts.DebugFunc != nil {
		debugFunc = opts.DebugFunc
	}

	fields := make([]string, 0)
	vars := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := ""

	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Insert(value.Elem().Interface(), opts)
	case reflect.Struct:
		if tab, ok := data.(Table); ok {
			table = tab.TableName()
		} else {
			table = value.Type().Name()
		}

		for i := 0; i < value.NumField(); i++ {
			if !value.Field(i).CanInterface() {
				continue
			}
			isTime := value.Field(i).Type().String() == "time.Time"

			tag := value.Type().Field(i).Tag.Get("xsql")
			if tag == "" || tag == "-" || tag == "_" {
				continue
			}
			fields = append(fields, tag)

			v := ""
			if placeholder == "?" {
				v = placeholder
			} else {
				v = fmt.Sprintf(placeholder, i)
			}
			if isTime {
				vars = append(vars, timeFunc(v))
			} else {
				vars = append(vars, v)
			}

			if isTime {
				ti := value.Field(i).Interface().(time.Time)
				bindArgs = append(bindArgs, ti.Format(timeLayout))
			} else {
				bindArgs = append(bindArgs, value.Field(i).Interface())
			}
		}
		break
	default:
		return nil, errors.New("sql: only for struct type")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES (%s)`, insertKey, table, columnQuotes+strings.Join(fields, columnQuotes+", "+columnQuotes)+columnQuotes, strings.Join(vars, `, `))

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Time:         time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if debugFunc != nil {
		debugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) BatchInsert(array interface{}, opts *Options) (sql.Result, error) {
	insertKey := "INSERT INTO"
	if opts.InsertKey != "" {
		insertKey = opts.InsertKey
	}
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeLayout := DefaultTimeLayout
	if opts.TimeLayout != "" {
		timeLayout = opts.TimeLayout
	}
	columnQuotes := "`"
	if opts.ColumnQuotes != "" {
		columnQuotes = opts.ColumnQuotes
	}
	var debugFunc DebugFunc
	if opts.DebugFunc != nil {
		debugFunc = opts.DebugFunc
	}

	fields := make([]string, 0)
	valueSql := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := ""

	// check
	value := reflect.ValueOf(array)
	switch value.Kind() {
	case reflect.Ptr:
		return t.BatchInsert(value.Elem().Interface(), opts)
	case reflect.Array, reflect.Slice:
		break
	default:
		return nil, errors.New("sql: only for struct array/slice type")
	}
	if value.Len() == 0 {
		return nil, errors.New("sql: array/slice length cannot be 0")
	}

	// fields
	switch value.Index(0).Kind() {
	case reflect.Struct:
		subValue := value.Index(0)

		if tab, ok := subValue.Interface().(Table); ok {
			table = tab.TableName()
		} else {
			table = subValue.Type().Name()
		}

		for i := 0; i < subValue.NumField(); i++ {
			if !subValue.Field(i).CanInterface() {
				continue
			}
			tag := subValue.Type().Field(i).Tag.Get("xsql")
			if tag == "" || tag == "-" || tag == "_" {
				continue
			}
			fields = append(fields, tag)
		}
		break
	default:
		return nil, errors.New("sql: only for struct array/slice type")
	}

	// values
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		ai := 0
		for r := 0; r < value.Len(); r++ {
			switch value.Index(r).Kind() {
			case reflect.Struct:
				subValue := value.Index(r)
				vars := make([]string, 0)
				for c := 0; c < subValue.NumField(); c++ {
					if !subValue.Field(c).CanInterface() {
						continue
					}

					tag := subValue.Type().Field(c).Tag.Get("xsql")
					if tag == "" || tag == "-" || tag == "_" {
						continue
					}

					if placeholder == "?" {
						vars = append(vars, placeholder)
					} else {
						vars = append(vars, fmt.Sprintf(placeholder, ai))
						ai += 1
					}

					// time特殊处理
					if subValue.Field(c).Type().String() == "time.Time" {
						ti := subValue.Field(c).Interface().(time.Time)
						bindArgs = append(bindArgs, ti.Format(timeLayout))
					} else {
						bindArgs = append(bindArgs, subValue.Field(c).Interface())
					}
				}
				valueSql = append(valueSql, fmt.Sprintf("(%s)", strings.Join(vars, `, `)))
				break
			default:
				return nil, errors.New("sql: only for struct array/slice type")
			}
		}
		break
	default:
		return nil, errors.New("sql: only for struct array/slice type")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES %s`, insertKey, table, columnQuotes+strings.Join(fields, columnQuotes+", "+columnQuotes)+columnQuotes, strings.Join(valueSql, ", "))

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Time:         time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if debugFunc != nil {
		debugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) Update(data interface{}, expr string, args []interface{}, opts *Options) (sql.Result, error) {
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeLayout := DefaultTimeLayout
	if opts.TimeLayout != "" {
		timeLayout = opts.TimeLayout
	}
	columnQuotes := "`"
	if opts.ColumnQuotes != "" {
		columnQuotes = opts.ColumnQuotes
	}
	var debugFunc DebugFunc
	if opts.DebugFunc != nil {
		debugFunc = opts.DebugFunc
	}

	set := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := ""

	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Update(value.Elem().Interface(), expr, args, opts)
	case reflect.Struct:
		if tab, ok := data.(Table); ok {
			table = tab.TableName()
		} else {
			table = value.Type().Name()
		}

		for i := 0; i < value.NumField(); i++ {
			if !value.Field(i).CanInterface() {
				continue
			}

			tag := value.Type().Field(i).Tag.Get("xsql")
			if tag == "" || tag == "-" || tag == "_" {
				continue
			}

			if placeholder == "?" {
				set = append(set, fmt.Sprintf("%s = %s", columnQuotes+tag+columnQuotes, placeholder))
			} else {
				set = append(set, fmt.Sprintf("%s = %s", columnQuotes+tag+columnQuotes, fmt.Sprintf(placeholder, i)))
			}

			// time特殊处理
			if value.Field(i).Type().String() == "time.Time" {
				ti := value.Field(i).Interface().(time.Time)
				bindArgs = append(bindArgs, ti.Format(timeLayout))
			} else {
				bindArgs = append(bindArgs, value.Field(i).Interface())
			}
		}
		break
	default:
		return nil, errors.New("sql: only for struct type")
	}

	where := ""
	if expr != "" {
		where = fmt.Sprintf(` WHERE %s`, expr)
		bindArgs = append(bindArgs, args...)
	}

	SQL := fmt.Sprintf(`UPDATE %s SET %s%s`, table, strings.Join(set, ", "), where)

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Time:         time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if debugFunc != nil {
		debugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) Exec(query string, args []interface{}, opts *Options) (sql.Result, error) {
	var debugFunc DebugFunc
	if opts.DebugFunc != nil {
		debugFunc = opts.DebugFunc
	}

	startTime := time.Now()
	res, err := t.Executor.Exec(query, args...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Time:         time.Now().Sub(startTime),
		SQL:          query,
		Bindings:     args,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if debugFunc != nil {
		debugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, err
}
