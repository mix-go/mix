package di

import (
	"github.com/mix-go/xdi"
	"net/http"
)

func init() {
	obj := xdi.Object{
		Name: "server",
		New: func() (i interface{}, e error) {
			return &http.Server{}, nil
		},
	}
	if err := xdi.Provide(&obj); err != nil {
		panic(err)
	}
}

func Server() (s *http.Server) {
	if err := xdi.Populate("server", &s); err != nil {
		panic(err)
	}
	return
}
