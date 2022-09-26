package myhttp

import (
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/gorilla/mux"
	"net/http"
)

var defaultBind = binding.New(nil)

func BindAndValidate(recvPointer interface{}, r *http.Request) error {
	return defaultBind.BindAndValidate(recvPointer, r, PathParams(mux.Vars(r)))
}

func Bind(recvPointer interface{}, r *http.Request) error {
	return defaultBind.Bind(recvPointer, r, PathParams(mux.Vars(r)))
}

func Validate(v interface{}) error {
	return defaultBind.Validate(v)
}

type PathParams map[string]string

func (p PathParams) Get(name string) (s string, ok bool) {
	s, ok = p[name]
	return
}
