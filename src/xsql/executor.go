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

func (t *executor) defaultOptions(opts *Options) *Options {
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
	newOpts := Options{
		InsertKey:    insertKey,
		Placeholder:  placeholder,
		ColumnQuotes: columnQuotes,
		TimeLayout:   timeLayout,
		TimeFunc:     timeFunc,
		DebugFunc:    debugFunc,
	}
	return &newOpts
}

func (t *executor) Insert(data interface{}, opts *Options) (sql.Result, error) {
	opts = t.defaultOptions(opts)

	fields := make([]string, 0)
	vars := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := ""

	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Insert(value.Elem().Interface(), opts)
	case reflect.Struct:
		if tab, ok := data.(Table); ok {
			table = tab.TableName()
		} else {
			table = value.Type().Name()
		}

		fields, vars, bindArgs = t.insertForeach(value, typ, opts)
		break
	default:
		return nil, errors.New("sql: only for struct type")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES (%s)`, opts.InsertKey, table, opts.ColumnQuotes+strings.Join(fields, opts.ColumnQuotes+", "+opts.ColumnQuotes)+opts.ColumnQuotes, strings.Join(vars, `, `))

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Duration:     time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if opts.DebugFunc != nil {
		opts.DebugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) BatchInsert(array interface{}, opts *Options) (sql.Result, error) {
	opts = t.defaultOptions(opts)

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
		subType := subValue.Type()

		if tab, ok := subValue.Interface().(Table); ok {
			table = tab.TableName()
		} else {
			table = subValue.Type().Name()
		}

		fields = t.batchInsertFields(subValue, subType)
		break
	default:
		return nil, errors.New("sql: only for struct array/slice type")
	}

	// values
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for r := 0; r < value.Len(); r++ {
			switch value.Index(r).Kind() {
			case reflect.Struct:
				subValue := value.Index(r)
				vars, b := t.batchInsertForeach(0, subValue, subValue.Type(), opts)
				bindArgs = append(bindArgs, b...)
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

	SQL := fmt.Sprintf(`%s %s (%s) VALUES %s`, opts.InsertKey, table, opts.ColumnQuotes+strings.Join(fields, opts.ColumnQuotes+", "+opts.ColumnQuotes)+opts.ColumnQuotes, strings.Join(valueSql, ", "))

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Duration:     time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if opts.DebugFunc != nil {
		opts.DebugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) Update(data interface{}, expr string, args []interface{}, opts *Options) (sql.Result, error) {
	opts = t.defaultOptions(opts)

	set := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := ""

	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Update(value.Elem().Interface(), expr, args, opts)
	case reflect.Struct:
		if tab, ok := data.(Table); ok {
			table = tab.TableName()
		} else {
			table = value.Type().Name()
		}

		set, bindArgs = t.updateForeach(value, typ, opts)
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
		Duration:     time.Now().Sub(startTime),
		SQL:          SQL,
		Bindings:     bindArgs,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if opts.DebugFunc != nil {
		opts.DebugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) Exec(query string, args []interface{}, opts *Options) (sql.Result, error) {
	opts = t.defaultOptions(opts)

	startTime := time.Now()
	res, err := t.Executor.Exec(query, args...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Duration:     time.Now().Sub(startTime),
		SQL:          query,
		Bindings:     args,
		RowsAffected: rowsAffected,
		Error:        err,
	}
	if opts.DebugFunc != nil {
		opts.DebugFunc(l)
	}
	if err != nil {
		return nil, err
	}

	return res, err
}

func (t *executor) insertForeach(value reflect.Value, typ reflect.Type, opts *Options) (fields, vars []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			f, v, b := t.insertForeach(fieldValue, fieldValue.Type(), opts)
			fields = append(fields, f...)
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}
		isTime := value.Field(n).Type().String() == "time.Time"

		tag := value.Type().Field(n).Tag.Get("xsql")
		if tag == "" || tag == "-" || tag == "_" {
			continue
		}
		fields = append(fields, tag)

		v := ""
		if opts.Placeholder == "?" {
			v = opts.Placeholder
		} else {
			v = fmt.Sprintf(opts.Placeholder, n)
		}
		if isTime {
			vars = append(vars, opts.TimeFunc(v))
		} else {
			vars = append(vars, v)
		}

		if isTime {
			ti := value.Field(n).Interface().(time.Time)
			bindArgs = append(bindArgs, ti.Format(opts.TimeLayout))
		} else {
			bindArgs = append(bindArgs, value.Field(n).Interface())
		}
	}
	return
}

func (t *executor) batchInsertFields(value reflect.Value, typ reflect.Type) (fields []string) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			f := t.batchInsertFields(fieldValue, fieldValue.Type())
			fields = append(fields, f...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get("xsql")
		if tag == "" || tag == "-" || tag == "_" {
			continue
		}

		fields = append(fields, tag)
	}
	return
}

func (t *executor) batchInsertForeach(ai int, value reflect.Value, typ reflect.Type, opts *Options) (vars []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			v, b := t.batchInsertForeach(ai+1000, fieldValue, fieldValue.Type(), opts)
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get("xsql")
		if tag == "" || tag == "-" || tag == "_" {
			continue
		}

		if opts.Placeholder == "?" {
			vars = append(vars, opts.Placeholder)
		} else {
			vars = append(vars, fmt.Sprintf(opts.Placeholder, ai))
			ai += 1
		}

		// time特殊处理
		if value.Field(n).Type().String() == "time.Time" {
			ti := value.Field(n).Interface().(time.Time)
			bindArgs = append(bindArgs, ti.Format(opts.TimeLayout))
		} else {
			bindArgs = append(bindArgs, value.Field(n).Interface())
		}
	}
	return
}

func (t *executor) updateForeach(value reflect.Value, typ reflect.Type, opts *Options) (set []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			s, b := t.updateForeach(fieldValue, fieldValue.Type(), opts)
			set = append(set, s...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get("xsql")
		if tag == "" || tag == "-" || tag == "_" {
			continue
		}

		if opts.Placeholder == "?" {
			set = append(set, fmt.Sprintf("%s = %s", opts.ColumnQuotes+tag+opts.ColumnQuotes, opts.Placeholder))
		} else {
			set = append(set, fmt.Sprintf("%s = %s", opts.ColumnQuotes+tag+opts.ColumnQuotes, fmt.Sprintf(opts.Placeholder, n)))
		}

		// time特殊处理
		if value.Field(n).Type().String() == "time.Time" {
			ti := value.Field(n).Interface().(time.Time)
			bindArgs = append(bindArgs, ti.Format(opts.TimeLayout))
		} else {
			bindArgs = append(bindArgs, value.Field(n).Interface())
		}
	}
	return
}
