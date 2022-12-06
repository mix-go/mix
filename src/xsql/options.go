package xsql

import "fmt"

// Options
// 默认为mysql模式
type Options struct {
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

// Oracle
// 使用oracle模式
func Oracle() Options {
	return Options{
		Placeholder:  `:%d`,
		ColumnQuotes: `"`,
		TimeFunc: func(placeholder string) string {
			return fmt.Sprintf("TO_TIMESTAMP(%s, 'SYYYY-MM-DD HH24:MI:SS:FF6')", placeholder)
		},
	}
}
