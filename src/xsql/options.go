package xsql

import "fmt"

var DefaultOptions = newDefaultOptions()

func newDefaultOptions() sqlOptions {
	return sqlOptions{
		Tag:          "xsql",
		InsertKey:    "INSERT INTO",
		Placeholder:  "?",
		ColumnQuotes: "`",
		TimeLayout:   "2006-01-02 15:04:05",
		TimeFunc: func(placeholder string) string {
			return placeholder
		},
		DebugFunc: nil,
	}
}

type sqlOptions struct {
	// 默认: xsql
	Tag string

	// 默认: INSERT INTO
	InsertKey string

	// 默认: ?
	// oracle 可配置为 :%d
	Placeholder string

	// 默认：`
	// oracle 可配置为 "
	ColumnQuotes string

	// 默认：== DefaultTimeLayout
	TimeLayout string

	// 默认：== DefaultTimeFunc
	// oracle 可修改这个闭包增加 TO_TIMESTAMP
	TimeFunc TimeFunc

	// 全局 debug SQL
	DebugFunc DebugFunc
}

func mergeOptions(opts []SqlOption) *sqlOptions {
	opt := DefaultOptions // copy
	for _, o := range opts {
		o.apply(&opt)
	}
	return &opt
}

type SqlOption interface {
	apply(*sqlOptions)
}

type funcSqlOption struct {
	f func(*sqlOptions)
}

func (fdo *funcSqlOption) apply(do *sqlOptions) {
	fdo.f(do)
}

func WithTag(tag string) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.Tag = tag
	}}
}

func WithInsertKey(insertKey string) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.InsertKey = insertKey
	}}
}

func WithPlaceholder(placeholder string) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.Placeholder = placeholder
	}}
}

func WithColumnQuotes(columnQuotes string) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.ColumnQuotes = columnQuotes
	}}
}

func WithTimeLayout(timeLayout string) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.TimeLayout = timeLayout
	}}
}

func WithTimeFunc(f TimeFunc) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.TimeFunc = f
	}}
}

func WithDebugFunc(f DebugFunc) SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.DebugFunc = f
	}}
}

func SwitchToOracle() SqlOption {
	return &funcSqlOption{func(opt *sqlOptions) {
		opt.Placeholder = `:%d`
		opt.ColumnQuotes = `"`
		opt.TimeFunc = func(placeholder string) string {
			return fmt.Sprintf("TO_TIMESTAMP(%s, 'SYYYY-MM-DD HH24:MI:SS:FF6')", placeholder)
		}
	}}
}
