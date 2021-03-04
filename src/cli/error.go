package cli

// NotFoundError 未找到错误
type NotFoundError struct {
	error
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}

// UnsupportError 不支持错误
type UnsupportError struct {
	error
}

// NewUnsupportError 创建不支持错误
func NewUnsupportError(err error) *UnsupportError {
	return &UnsupportError{err}
}
