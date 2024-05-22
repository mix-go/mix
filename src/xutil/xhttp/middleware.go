package xhttp

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(*Request, *RequestOptions) (*Response, error)
