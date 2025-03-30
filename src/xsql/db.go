package xsql

import (
	"database/sql"
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
)

type TimeFunc func(placeholder string) string

type DB struct {
	Options  *sqlOptions
	raw      *sql.DB
	executor executor
	query    query
}

func New(db *sql.DB, opts ...SqlOption) *DB {
	return &DB{
		Options: mergeOptions(opts),
		raw:     db,
		executor: executor{
			Executor: db,
		},
		query: query{
			Query: db,
		},
	}
}

func (t *DB) mergeOptions(opts []SqlOption) *sqlOptions {
	cp := *t.Options // copy
	for _, o := range opts {
		o.apply(&cp)
	}
	return &cp
}

func (t *DB) Insert(data interface{}, opts ...SqlOption) (sql.Result, error) {
	return t.executor.Insert(data, t.mergeOptions(opts))
}

func (t *DB) BatchInsert(data interface{}, opts ...SqlOption) (sql.Result, error) {
	return t.executor.BatchInsert(data, t.mergeOptions(opts))
}

func (t *DB) Update(data interface{}, expr string, args ...interface{}) (sql.Result, error) {
	return t.executor.Update(data, expr, args, t.Options)
}

func (t *DB) Model(s interface{}) *ModelExecutor {
	return t.executor.model(s, t.Options)
}

func (t *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.executor.Exec(query, args, t.Options)
}

func (t *DB) Begin() (*Tx, error) {
	tx, err := t.raw.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{
		raw: tx,
		DB: &DB{
			Options: t.Options,
			executor: executor{
				Executor: tx,
			},
			query: query{
				Query: tx,
			},
		},
	}, nil
}

func (t *DB) Query(query string, args ...interface{}) ([]*Row, error) {
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return nil, err
	}
	r, err := f.Rows()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (t *DB) QueryFirst(query string, args ...interface{}) (*Row, error) {
	rows, err := t.Query(query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}
	return rows[0], nil
}

func (t *DB) Find(i interface{}, query string, args ...interface{}) error {
	query = tableReplace(i, query, t.Options)
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return err
	}
	if err := f.Find(i); err != nil {
		return err
	}
	return nil
}

func (t *DB) First(i interface{}, query string, args ...interface{}) error {
	query = tableReplace(i, query, t.Options)
	f, err := t.query.Fetch(query, args, t.Options)
	if err != nil {
		return err
	}
	if err := f.First(i); err != nil {
		return err
	}
	return nil
}

func tableReplace(i interface{}, query string, opts *sqlOptions) string {
	var table string

	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().IsValid() {
			// *Test > Test
			if value.Elem().Kind() == reflect.Struct {
				// 先尝试*Test能不能找到
				if tab, ok := value.Interface().(Table); ok {
					table = tab.TableName()
					break
				}
			}
			// **Test > *Test
			return tableReplace(value.Elem().Interface(), query, opts)
		}
		if tab, ok := value.Interface().(Table); ok {
			table = tab.TableName()
			break
		}
		table = getTypeName(i)
	case reflect.Struct:
		if tab, ok := value.Interface().(Table); ok {
			table = tab.TableName()
			break
		}
		// 也去尝试*Test能不能找到
		valuePtr := reflect.New(value.Type())
		if tab, ok := valuePtr.Interface().(Table); ok {
			table = tab.TableName()
			break
		}
		table = getTypeName(i)
	case reflect.Array, reflect.Slice:
		elemType := value.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			// 创建该类型的新实例
			// Test
			elemValue := reflect.New(elemType)
			elemInstance := elemValue.Interface()
			if tab, ok := elemInstance.(Table); ok {
				table = tab.TableName()
				break
			}
			// Test > *Test
			elemPtrInstance := elemValue.Addr().Interface()
			if tab, ok := elemPtrInstance.(Table); ok {
				table = tab.TableName()
				break
			}
			table = getTypeName(elemInstance)
		} else {
			return query // 如果元素不是结构体或其指针，返回原始查询
		}
	default:
		return query // 如果不是结构体、数组或切片，返回原始查询
	}

	return strings.Replace(query, opts.TableKey, table, 1)
}

func getTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

var ProtoMarshalOptions = protojson.MarshalOptions{
	UseProtoNames:  true,
	UseEnumNumbers: true,
}

func marshal(v any) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return ProtoMarshalOptions.Marshal(m)
	} else {
		return json.Marshal(v)
	}
}

var ProtoUnmarshalOptions = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func unmarshal(b []byte, v any) error {
	if m, ok := v.(proto.Message); ok {
		return ProtoUnmarshalOptions.Unmarshal(b, m)
	} else {
		return json.Unmarshal(b, v)
	}
}
