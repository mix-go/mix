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
	DB *sql.DB
}

func (t *executor) Insert(data interface{}, opts *Options) (sql.Result, *Log, error) {
	insertKey := "INSERT INTO"
	if opts.InsertKey != "" {
		insertKey = opts.InsertKey
	}
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeParseLayout := DefaultTimeParseLayout
	if opts.TimeParseLayout != "" {
		timeParseLayout = opts.TimeParseLayout
	}
	quoteSymbol := "`"
	if opts.QuoteSymbol != "" {
		quoteSymbol = opts.QuoteSymbol
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

			tag := value.Type().Field(i).Tag.Get("xsql")
			if tag == "" || tag == "-" || tag == "_" {
				continue
			}

			fields = append(fields, tag)

			if placeholder == "?" {
				vars = append(vars, placeholder)
			} else {
				vars = append(vars, fmt.Sprintf(placeholder, i))
			}

			// time特殊处理
			if value.Field(i).Type().String() == "time.Time" {
				ti := value.Field(i).Interface().(time.Time)
				bindArgs = append(bindArgs, ti.Format(timeParseLayout))
			} else {
				bindArgs = append(bindArgs, value.Field(i).Interface())
			}
		}
		break
	default:
		return nil, nil, errors.New("only for struct type")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES (%s)`, insertKey, table, quoteSymbol+strings.Join(fields, quoteSymbol+", "+quoteSymbol)+quoteSymbol, strings.Join(vars, `, `))

	startTime := time.Now()
	res, err := t.DB.Exec(SQL, bindArgs...)
	l := &Log{
		SQL:  SQL,
		Args: bindArgs,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}

	return res, l, nil
}

func (t *executor) BatchInsert(array interface{}, opts *Options) (sql.Result, *Log, error) {
	insertKey := "INSERT INTO"
	if opts.InsertKey != "" {
		insertKey = opts.InsertKey
	}
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeParseLayout := DefaultTimeParseLayout
	if opts.TimeParseLayout != "" {
		timeParseLayout = opts.TimeParseLayout
	}
	quoteSymbol := "`"
	if opts.QuoteSymbol != "" {
		quoteSymbol = opts.QuoteSymbol
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
		return nil, nil, errors.New("only for struct array/slice type")
	}
	if value.Len() == 0 {
		return nil, nil, errors.New("array/slice length cannot be 0")
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
		return nil, nil, errors.New("only for struct array/slice type")
	}

	// values
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			switch value.Index(i).Kind() {
			case reflect.Struct:
				subValue := value.Index(i)
				vars := make([]string, 0)
				for i := 0; i < subValue.NumField(); i++ {
					if !subValue.Field(i).CanInterface() {
						continue
					}

					tag := subValue.Type().Field(i).Tag.Get("xsql")
					if tag == "" || tag == "-" || tag == "_" {
						continue
					}

					if placeholder == "?" {
						vars = append(vars, placeholder)
					} else {
						vars = append(vars, fmt.Sprintf(placeholder, i))
					}

					// time特殊处理
					if subValue.Field(i).Type().String() == "time.Time" {
						ti := subValue.Field(i).Interface().(time.Time)
						bindArgs = append(bindArgs, ti.Format(timeParseLayout))
					} else {
						bindArgs = append(bindArgs, subValue.Field(i).Interface())
					}
				}
				valueSql = append(valueSql, fmt.Sprintf("(%s)", strings.Join(vars, `, `)))
				break
			default:
				return nil, nil, errors.New("only for struct array/slice type")
			}
		}
		break
	default:
		return nil, nil, errors.New("only for struct array/slice type")
	}

	SQL := fmt.Sprintf(`%s %s (%s) VALUES %s`, insertKey, table, quoteSymbol+strings.Join(fields, quoteSymbol+", "+quoteSymbol)+quoteSymbol, strings.Join(valueSql, ", "))

	startTime := time.Now()
	res, err := t.DB.Exec(SQL, bindArgs...)
	l := &Log{
		SQL:  SQL,
		Args: bindArgs,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}

	return res, l, nil
}

func (t *executor) Update(data interface{}, expr string, args []interface{}, opts *Options) (sql.Result, *Log, error) {
	placeholder := "?"
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}
	timeParseLayout := DefaultTimeParseLayout
	if opts.TimeParseLayout != "" {
		timeParseLayout = opts.TimeParseLayout
	}
	quoteSymbol := "`"
	if opts.QuoteSymbol != "" {
		quoteSymbol = opts.QuoteSymbol
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
				set = append(set, fmt.Sprintf("%s = %s", quoteSymbol+tag+quoteSymbol, placeholder))
			} else {
				set = append(set, fmt.Sprintf("%s = %s", quoteSymbol+tag+quoteSymbol, fmt.Sprintf(placeholder, i)))
			}

			// time特殊处理
			if value.Field(i).Type().String() == "time.Time" {
				ti := value.Field(i).Interface().(time.Time)
				bindArgs = append(bindArgs, ti.Format(timeParseLayout))
			} else {
				bindArgs = append(bindArgs, value.Field(i).Interface())
			}
		}
		break
	default:
		return nil, nil, errors.New("only for struct type")
	}

	where := ""
	if expr != "" {
		where = fmt.Sprintf(` WHERE %s`, expr)
		bindArgs = append(bindArgs, args...)
	}

	SQL := fmt.Sprintf(`UPDATE %s SET %s%s`, table, strings.Join(set, ", "), where)

	startTime := time.Now()
	res, err := t.DB.Exec(SQL, bindArgs...)
	l := &Log{
		SQL:  SQL,
		Args: bindArgs,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}

	return res, l, nil
}
