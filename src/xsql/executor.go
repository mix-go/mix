package xsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ora "github.com/sijms/go-ora/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"slices"
	"strings"
	"time"
)

type executor struct {
	Executor
}

func (t *executor) Insert(ctx context.Context, data interface{}, opts *sqlOptions) (sql.Result, error) {
	var err error
	fields := make([]string, 0)
	vars := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Insert(ctx, value.Elem().Interface(), opts)
	case reflect.Struct:
		fields, vars, bindArgs, err = t.foreachInsert(value, typ, opts)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("xsql: only support struct")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES (%s)`, opts.InsertKey, opts.TableKey, opts.ColumnQuotes+strings.Join(fields, opts.ColumnQuotes+", "+opts.ColumnQuotes)+opts.ColumnQuotes, strings.Join(vars, `, `))
	SQL = tableReplace(data, SQL, opts)

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Context:      ctx,
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

func (t *executor) BatchInsert(ctx context.Context, array interface{}, opts *sqlOptions) (sql.Result, error) {
	fields := make([]string, 0)
	valueSql := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	// check
	value := reflect.ValueOf(array)
	switch value.Kind() {
	case reflect.Ptr:
		return t.BatchInsert(ctx, value.Elem().Interface(), opts)
	case reflect.Array, reflect.Slice:
		break
	default:
		return nil, errors.New("xsql: only support array, slice")
	}
	if value.Len() == 0 {
		return nil, errors.New("xsql: array, slice length cannot be 0")
	}

	// fields
	switch value.Index(0).Kind() {
	case reflect.Struct:
		subValue := value.Index(0)
		subType := subValue.Type()
		fields = t.foreachBatchInsertFields(subValue, subType, opts)
		break
	default:
		return nil, errors.New("xsql: only support array, slice")
	}

	// values
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for r := 0; r < value.Len(); r++ {
			switch value.Index(r).Kind() {
			case reflect.Struct:
				subValue := value.Index(r)
				vars, b, err := t.foreachBatchInsertValues(0, subValue, subValue.Type(), opts)
				if err != nil {
					return nil, err
				}
				bindArgs = append(bindArgs, b...)
				valueSql = append(valueSql, fmt.Sprintf("(%s)", strings.Join(vars, `, `)))
				break
			default:
				return nil, errors.New("xsql: only support array, slice")
			}
		}
		break
	default:
		return nil, errors.New("xsql: only support array, slice")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES %s`, opts.InsertKey, opts.TableKey, opts.ColumnQuotes+strings.Join(fields, opts.ColumnQuotes+", "+opts.ColumnQuotes)+opts.ColumnQuotes, strings.Join(valueSql, ", "))
	SQL = tableReplace(array, SQL, opts)

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Context:      ctx,
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
	table := tableReplace(s, opts.TableKey, opts)
	return &ModelExecutor{
		Executor:  t.Executor,
		Options:   opts,
		TableName: table,
	}
}

func (t *executor) Update(ctx context.Context, data interface{}, expr string, args []interface{}, opts *sqlOptions) (sql.Result, error) {
	var err error
	set := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	switch value.Kind() {
	case reflect.Ptr:
		return t.Update(ctx, value.Elem().Interface(), expr, args, opts)
	case reflect.Struct:
		set, bindArgs, err = t.foreachUpdate(value, typ, opts)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("xsql: only support struct")
	}

	where := ""
	if expr != "" {
		where = fmt.Sprintf(` WHERE %s`, expr)
		bindArgs = append(bindArgs, args...)
	}

	SQL := fmt.Sprintf(`UPDATE %s SET %s%s`, opts.TableKey, strings.Join(set, ", "), where)
	SQL = tableReplace(data, SQL, opts)

	startTime := time.Now()
	res, err := t.Executor.Exec(SQL, bindArgs...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Context:      ctx,
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

func (t *executor) Exec(ctx context.Context, query string, args []interface{}, opts *sqlOptions) (sql.Result, error) {
	startTime := time.Now()
	res, err := t.Executor.Exec(query, args...)
	var rowsAffected int64
	if res != nil {
		rowsAffected, _ = res.RowsAffected()
	}
	l := &Log{
		Context:      ctx,
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

func isTime(typ string) bool {
	switch typ {
	case "time.Time", "go_ora.TimeStamp", "*timestamppb.Timestamp":
		return true
	default:
		return false
	}
}

func formatTime(typ string, v interface{}, opts *sqlOptions) string {
	switch typ {
	case "time.Time":
		return v.(time.Time).Format(opts.TimeLayout)
	case "go_ora.TimeStamp":
		return time.Time(v.(ora.TimeStamp)).Format(opts.TimeLayout)
	case "*timestamppb.Timestamp":
		return v.(*timestamppb.Timestamp).AsTime().Format(opts.TimeLayout)
	default:
		return ""
	}
}

func (t *executor) foreachInsert(value reflect.Value, typ reflect.Type, opts *sqlOptions) (fields, vars []string, bindArgs []interface{}, err error) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)

		// Embedded Structs
		if fieldStruct.Anonymous {
			f, v, b, e := t.foreachInsert(fieldValue, fieldValue.Type(), opts)
			if e != nil {
				return nil, nil, nil, e
			}
			fields = append(fields, f...)
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !fieldValue.CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		fields = append(fields, tag)

		vTyp := fieldValue.Type().String()
		isTime := isTime(vTyp)

		var v string
		if opts.Placeholder == "?" {
			v = opts.Placeholder
		} else {
			v = fmt.Sprintf(opts.Placeholder, n)
		}
		if isTime {
			v = opts.TimeFunc(v)
		}
		vars = append(vars, v)

		var a interface{}
		if isTime {
			a = formatTime(vTyp, fieldValue.Interface(), opts)
		} else {
			// 非标量用JSON序列化处理
			if slices.Contains([]reflect.Kind{reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map}, fieldValue.Kind()) {
				b, e := marshal(fieldValue.Interface())
				if e != nil {
					return nil, nil, nil, fmt.Errorf("json unmarshal error %s for field %s", e, tag)
				}
				a = string(b)
			} else {
				a = fieldValue.Interface()
			}
		}
		bindArgs = append(bindArgs, a)
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

		if !fieldValue.CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		fields = append(fields, tag)
	}
	return
}

func (t *executor) foreachBatchInsertValues(ai int, value reflect.Value, typ reflect.Type, opts *sqlOptions) (vars []string, bindArgs []interface{}, err error) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)

		// Embedded Structs
		if fieldStruct.Anonymous {
			v, b, e := t.foreachBatchInsertValues(ai+1000, fieldValue, fieldValue.Type(), opts)
			if e != nil {
				return nil, nil, e
			}
			vars = append(vars, v...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !fieldValue.CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		vTyp := fieldValue.Type().String()
		isTime := isTime(vTyp)

		var v string
		if opts.Placeholder == "?" {
			v = opts.Placeholder
		} else {
			v = fmt.Sprintf(opts.Placeholder, ai)
			ai += 1
		}
		if isTime {
			v = opts.TimeFunc(v)
		}
		vars = append(vars, v)

		var a interface{}
		if isTime {
			a = formatTime(vTyp, fieldValue.Interface(), opts)
		} else {
			// 非标量用JSON序列化处理
			if slices.Contains([]reflect.Kind{reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map}, fieldValue.Kind()) {
				b, e := marshal(fieldValue.Interface())
				if e != nil {
					return nil, nil, fmt.Errorf("json unmarshal error %s for field %s", e, tag)
				}
				a = string(b)
			} else {
				a = fieldValue.Interface()
			}
		}
		bindArgs = append(bindArgs, a)
	}
	return
}

func (t *executor) foreachUpdate(value reflect.Value, typ reflect.Type, opts *sqlOptions) (set []string, bindArgs []interface{}, err error) {
	for n := 0; n < value.NumField(); n++ {
		fieldValue := value.Field(n)
		fieldStruct := typ.Field(n)

		// Embedded Structs
		if fieldStruct.Anonymous {
			s, b, e := t.foreachUpdate(fieldValue, fieldValue.Type(), opts)
			if e != nil {
				return nil, nil, e
			}
			set = append(set, s...)
			bindArgs = append(bindArgs, b...)
			continue
		}

		if !fieldValue.CanInterface() {
			continue
		}

		tag := value.Type().Field(n).Tag.Get(opts.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		vTyp := fieldValue.Type().String()
		isTime := isTime(vTyp)

		var v string
		if opts.Placeholder == "?" {
			v = opts.Placeholder
		} else {
			v = fmt.Sprintf(opts.Placeholder, n)
		}
		if isTime {
			v = opts.TimeFunc(v)
		}
		set = append(set, fmt.Sprintf("%s = %s", opts.ColumnQuotes+tag+opts.ColumnQuotes, v))

		var a interface{}
		if isTime {
			a = formatTime(vTyp, fieldValue.Interface(), opts)
		} else {
			// 非标量用JSON序列化处理
			if slices.Contains([]reflect.Kind{reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map}, fieldValue.Kind()) {
				b, e := marshal(fieldValue.Interface())
				if e != nil {
					return nil, nil, fmt.Errorf("json unmarshal error %s for field %s", e, tag)
				}
				a = string(b)
			} else {
				a = fieldValue.Interface()
			}
		}
		bindArgs = append(bindArgs, a)
	}
	return
}
