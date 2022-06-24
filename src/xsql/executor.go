package xsql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Executor struct {
	DB *sql.DB
}

func (t *Executor) Insert(table string, data interface{}, opts ...Option) (sql.Result, *Log, error) {
	config := &Config{}

	for _, opt := range opts {
		if opt != nil {
			opt.Apply(config)
		}
	}

	insertKey := "INSERT INTO"
	if config.InsertKey != "" {
		insertKey = config.InsertKey
	}
	placeholder := "?"
	if config.Placeholder != "" {
		placeholder = config.Placeholder
	}

	fields := make([]string, 0)
	vars := make([]string, 0)
	args := make([]interface{}, 0)

	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Struct:
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
			args = append(args, value.Field(i).Interface())
		}
		break
	default:
		return nil, nil, errors.New("only for struct type")
	}

	s := fmt.Sprintf(`%s %s ("%s") VALUES (%s)`, insertKey, table, strings.Join(fields, `", "`), strings.Join(vars, `, `))

	startTime := time.Now()
	res, err := t.DB.Exec(s, args...)
	l := &Log{
		SQL:  s,
		Args: args,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}

	return res, l, nil
}

func (t *Executor) BatchInsert(table string, array interface{}, opts ...Option) (sql.Result, *Log, error) {
	config := &Config{}

	for _, opt := range opts {
		if opt != nil {
			opt.Apply(config)
		}
	}

	insertKey := "INSERT INTO"
	if config.InsertKey != "" {
		insertKey = config.InsertKey
	}
	placeholder := "?"
	if config.Placeholder != "" {
		placeholder = config.Placeholder
	}

	fields := make([]string, 0)
	valueSql := make([]string, 0)
	args := make([]interface{}, 0)

	value := reflect.ValueOf(array)
	if value.Len() == 0 {
		return nil, nil, errors.New("array/slice length cannot be 0")
	}
	switch value.Index(0).Kind() {
	case reflect.Struct:
		subValue := value.Index(0)
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
					args = append(args, subValue.Field(i).Interface())
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

	s := fmt.Sprintf(`%s %s ("%s") VALUES %s`, insertKey, table, strings.Join(fields, `", "`), strings.Join(valueSql, `", "`))

	startTime := time.Now()
	res, err := t.DB.Exec(s, args...)
	l := &Log{
		SQL:  s,
		Args: args,
		Time: time.Now().Sub(startTime),
	}
	if err != nil {
		return nil, l, err
	}

	return res, l, nil
}
