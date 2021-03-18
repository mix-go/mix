package xcli

// NotFoundError 未找到错误
type NotFoundError struct {
	error
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}

// UnsupportedError 不支持错误
type UnsupportedError struct {
	error
}

// NewUnsupportedError 创建不支持错误
func NewUnsupportedError(err error) *UnsupportedError {
	return &UnsupportedError{err}
}
