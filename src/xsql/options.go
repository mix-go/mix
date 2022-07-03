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

	// 默认：== DefaultTimeParseLayout
	TimeParseLayout string
}
