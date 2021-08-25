package diyhttprouter

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	StatusBadRequest = 400 // RFC 7231, 6.5.1

)

type Handle func(http.ResponseWriter, *http.Request)

type router struct {
	trees map[string]*node
}

func New() *router {
	return &router{}
}

func (r *router) GET(path string, handle Handle) {
	r.Handle("GET", path, handle)
}

func (r *router) Handle(method, path string, handle Handle) {
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/'")
	}
	if r.trees == nil {
		r.trees = make(map[string]*node)
	}
	if r.trees[method] == nil {
		r.trees[method] = &node{
			path:   path,
			handle: handle,
		}
		return
	}
	r.trees[method].insertNode(path, handle)
}

func (r *router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "*" {
		if req.ProtoAtLeast(1, 1) {
			rw.Header().Set("Connection", "close")
		}
		rw.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := r.Handler(req)
	h(rw, req)

}

func (r *router) Handler(req *http.Request) (Handle, error) {
	if r.trees == nil {
		return nil, errors.New("no route registed")
	}
	if r.trees[req.Method] == nil {
		return nil, errors.New(fmt.Sprintf("no route for method %s registed", req.Method))
	}
	return r.trees[req.Method].getValue(req.URL.Path)
}
