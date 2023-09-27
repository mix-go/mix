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
	// Default: xsql
	Tag string

	// Default: INSERT INTO
	InsertKey string

	// Default: ?
	// For oracle, can be configured as :%d
	Placeholder string

	// Default: `
	// For oracle, can be configured as "
	ColumnQuotes string

	// Default: == DefaultTimeLayout
	TimeLayout string

	// Default: == DefaultTimeFunc
	// For oracle, this closure can be modified to add TO_TIMESTAMP
	TimeFunc TimeFunc

	// Global debug SQL
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
