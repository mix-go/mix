package xsql

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

type ModelExecutor struct {
	Executor
	Options *sqlOptions

	TableName string

	Error        error
	RowsAffected int64
}

func (t *ModelExecutor) Update(data map[string]interface{}, expr string, args ...interface{}) *ModelExecutor {
	return t.getExecResult(t.update(data, expr, args...))
}

func (t *ModelExecutor) getExecResult(r sql.Result, err error) *ModelExecutor {
	if err != nil {
		t.Error = err
		return t
	} else {
		t.Error = nil
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		t.Error = err
		return t
	} else {
		t.Error = nil
	}
	t.RowsAffected = rowsAffected
	return t
}

func (t *ModelExecutor) update(data map[string]interface{}, expr string, args ...interface{}) (sql.Result, error) {
	set := make([]string, 0)
	bindArgs := make([]interface{}, 0)

	table := t.TableName
	opts := t.Options

	n := 0
	for key, val := range data {
		value := reflect.ValueOf(val)
		vTyp := value.Type().String()
		isTime := isTime(vTyp)

		var v string
		if opts.Placeholder == "?" {
			v = opts.Placeholder
		} else {
			v = fmt.Sprintf(opts.Placeholder, n)
			n++
		}
		if isTime {
			v = opts.TimeFunc(v)
		}

		var a interface{}
		if isTime {
			a = formatTime(vTyp, value.Interface(), opts)
		} else {
			// 非标量用JSON序列化处理
			if slices.Contains([]reflect.Kind{reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map}, value.Kind()) {
				b, e := marshal(value.Interface())
				if e != nil {
					return nil, fmt.Errorf("json unmarshal error %s for field %s", e, key)
				}
				a = string(b)
			} else {
				a = value.Interface()
			}
		}

		set = append(set, fmt.Sprintf("`%s` = %s", key, v))
		bindArgs = append(bindArgs, a)
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

func (t *ModelExecutor) Delete(expr string, args ...interface{}) *ModelExecutor {
	return t.getExecResult(t.delete(expr, args...))
}

func (t *ModelExecutor) delete(expr string, args ...interface{}) (sql.Result, error) {
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
