package xsql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type executor struct {
	Executor
}

type ModelExecutor struct {
	Executor
	Options   *sqlOptions
	TableName string
}

func (t *executor) Insert(data interface{}, opts *sqlOptions) (sql.Result, error) {
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

		fields, vars, bindArgs = t.foreachInsert(value, typ, opts)
		break
	default:
		return nil, errors.New("xsql: only for struct type")
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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) BatchInsert(array interface{}, opts *sqlOptions) (sql.Result, error) {
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
		return nil, errors.New("xsql: only for struct array/slice type")
	}
	if value.Len() == 0 {
		return nil, errors.New("xsql: array/slice length cannot be 0")
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

		fields = t.foreachBatchInsertFields(subValue, subType, opts)
		break
	default:
		return nil, errors.New("xsql: only for struct array/slice type")
	}

	// values
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for r := 0; r < value.Len(); r++ {
			switch value.Index(r).Kind() {
			case reflect.Struct:
				subValue := value.Index(r)
				vars, b := t.foreachBatchInsertValues(0, subValue, subValue.Type(), opts)
				bindArgs = append(bindArgs, b...)
				valueSql = append(valueSql, fmt.Sprintf("(%s)", strings.Join(vars, `, `)))
				break
			default:
				return nil, errors.New("xsql: only for struct array/slice type")
			}
		}
		break
	default:
		return nil, errors.New("xsql: only for struct array/slice type")
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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) model(s interface{}, opts *sqlOptions) *ModelExecutor {
	var table string
	value := reflect.ValueOf(s)
	switch value.Kind() {
	case reflect.Ptr:
		return t.model(value.Elem().Interface(), opts)
	case reflect.Struct:
		if tab, ok := s.(Table); ok {
			table = tab.TableName()
		} else {
			table = value.Type().Name()
		}
		break
	}
	return &ModelExecutor{
		Executor:  t.Executor,
		Options:   opts,
		TableName: table,
	}
}

func (t *executor) Update(data interface{}, expr string, args []interface{}, opts *sqlOptions) (sql.Result, error) {
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

		set, bindArgs = t.foreachUpdate(value, typ, opts)
		break
	default:
		return nil, errors.New("xsql: only for struct type")
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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *ModelExecutor) Update(data map[string]interface{}, expr string, args ...interface{}) (sql.Result, error) {
	set := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := t.TableName
	opts := t.Options

	for k, v := range data {
		set = append(set, fmt.Sprintf("`%s` = ?", k))
		bindArgs = append(bindArgs, v)
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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *ModelExecutor) Delete(expr string, args ...interface{}) (sql.Result, error) {
	bindArgs := make([]interface{}, 0)

	table := t.TableName
	opts := t.Options

	where := ""
	if expr != "" {
		where = fmt.Sprintf(` WHERE %s`, expr)
		bindArgs = append(bindArgs, args...)
	}

	SQL := fmt.Sprintf(`DELETE FROM %s%s`, table, where)

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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *executor) Exec(query string, args []interface{}, opts *sqlOptions) (sql.Result, error) {
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
	opts.doDebug(l)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (t *executor) foreachInsert(value reflect.Value, typ reflect.Type, opts *sqlOptions) (fields, vars []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			f, v, b := t.foreachInsert(fieldValue, fieldValue.Type(), opts)
			fields = append(fields, f...)
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}
		isTime := value.Field(n).Type().String() == "time.Time"

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
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

func (t *executor) foreachBatchInsertFields(value reflect.Value, typ reflect.Type, opts *sqlOptions) (fields []string) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			f := t.foreachBatchInsertFields(fieldValue, fieldValue.Type(), opts)
			fields = append(fields, f...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
		if tag == "" || tag == "-" || tag == "_" {
			continue
		}

		fields = append(fields, tag)
	}
	return
}

func (t *executor) foreachBatchInsertValues(ai int, value reflect.Value, typ reflect.Type, opts *sqlOptions) (vars []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			v, b := t.foreachBatchInsertValues(ai+1000, fieldValue, fieldValue.Type(), opts)
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
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

func (t *executor) foreachUpdate(value reflect.Value, typ reflect.Type, opts *sqlOptions) (set []string, bindArgs []interface{}) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)
		if fieldStruct.Anonymous {
			s, b := t.foreachUpdate(fieldValue, fieldValue.Type(), opts)
			set = append(set, s...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !value.Field(n).CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
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
