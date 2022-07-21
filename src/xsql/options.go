package xsql

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
	// oracle 可修改这个闭包
	TimeFunc TimeFunc

	// 全局 debug SQL
	DebugFunc DebugFunc
}
