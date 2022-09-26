package myhttpmid

import (
	"context"
	"fmt"
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/fitan/mykit/myhttp"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"sync"
)

type debugSwitch struct {
	l    []msg
	m    map[string]bool
	lock sync.RWMutex
}

func NewDebugSwitch() *debugSwitch {
	return &debugSwitch{
		l: make([]msg, 0),
		m: make(map[string]bool, 0),
	}
}

type msg struct {
	Annotation string
	HttpPath   string
	Method     string
	Enable     bool
}

type Request struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Enable bool   `json:"enable"`
}

func (d *debugSwitch) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c := mux.CurrentRoute(r)
		path, _ := c.GetPathTemplate()
		methods, _ := c.GetMethods()
		for _, m := range methods {
			has, _ := d.Debug(path, m)
			if !has {
				continue
			}
			r = r.WithContext(context.WithValue(r.Context(), ContextKeyDebugEnable, has))
			next.ServeHTTP(w, r)
			break
		}
	}
	return http.HandlerFunc(fn)
}

func (d *debugSwitch) Handlers(prefix string, mux *mux.Router) {
	r := mux.PathPrefix(prefix).Subrouter()

	r.HandleFunc("", func(writer http.ResponseWriter, request *http.Request) {
		req := Request{}
		err := binding.New(nil).Bind(&req, request, nil)
		if err != nil {
			myhttp.ResponseJsonEncode(writer, err.Error())
			return
		}

		err = d.SetDebug(req.Path, req.Method, req.Enable)
		if err != nil {
			myhttp.ResponseJsonEncode(writer, err.Error())
			return
		}

		myhttp.ResponseJsonEncode(writer, "ok")
	}).Methods(http.MethodPut)

	r.HandleFunc("/reset", func(writer http.ResponseWriter, request *http.Request) {
		d.ResetDebug()
		myhttp.ResponseJsonEncode(writer, "ok")
	}).Methods(http.MethodPost)

	r.HandleFunc("", func(writer http.ResponseWriter, request *http.Request) {
		myhttp.ResponseJsonEncode(writer, d.List())
	}).Methods(http.MethodGet)
}

func (d *debugSwitch) Register(annotation, httpPath, method string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.m[path.Join(method, httpPath)] = false
	d.l = append(d.l, msg{
		Annotation: annotation,
		HttpPath:   httpPath,
		Method:     method,
	})
}

func (d *debugSwitch) SetDebug(httpPath, method string, b bool) (err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.m[path.Join(method, httpPath)]
	if ok {
		d.m[path.Join(method, httpPath)] = b
		return nil
	}

	return fmt.Errorf("not locate path and method")
}

func (d *debugSwitch) ResetDebug() {
	d.lock.Lock()
	defer d.lock.Unlock()
	for index, _ := range d.m {
		d.m[index] = false
	}
	return
}

func (d *debugSwitch) Debug(httpPath, method string) (bool, error) {
	has, ok := d.m[path.Join(method, httpPath)]
	if ok {
		return has, nil
	}
	return false, fmt.Errorf("not locate path and method")

}

func (d *debugSwitch) List() []msg {
	for i, _ := range d.l {
		d.l[i].Enable = d.m[d.l[i].HttpPath+d.l[i].Method]
	}
	return d.l
}
