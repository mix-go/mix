package xhttp

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(*XRequest, *RequestOptions) (*XResponse, error)
