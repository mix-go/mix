package xhttp

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(*Request) (*Response, error)
