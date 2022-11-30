package msgo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandleFunc func(ctx *Context)

type router struct {
	groups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{
		groupName:        name,
		handlerMap:       make(map[string]map[string]HandleFunc),
		handlerMethodMap: make(map[string][]string),
	}
	r.groups = append(r.groups, g)
	return g
}

func (r *routerGroup) handle(name string, method string, handleFunc HandleFunc) {
	_, ok := r.handlerMap[name]
	if !ok {
		r.handlerMap[name] = make(map[string]HandleFunc)
	}
	r.handlerMap[name][method] = handleFunc
	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], name)
}

func (r *routerGroup) Any(name string, handleFunc HandleFunc) {

	r.handle(name, ANY, handleFunc)
}

func (r *routerGroup) Handle(name string, method string, handleFunc HandleFunc) {

	r.handle(name, method, handleFunc)
}

func (r *routerGroup) Get(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodGet, handleFunc)
}
func (r *routerGroup) Post(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPost, handleFunc)
}
func (r *routerGroup) Patch(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPatch, handleFunc)
}
func (r *routerGroup) Delete(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodDelete, handleFunc)
}
func (r *routerGroup) Put(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPut, handleFunc)
}
func (r *routerGroup) Options(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodOptions, handleFunc)
}
func (r *routerGroup) Head(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodHead, handleFunc)
}

type routerGroup struct {
	groupName        string
	handlerMap       map[string]map[string]HandleFunc
	handlerMethodMap map[string][]string
}

type Engine struct {
	*router
}

func New() *Engine {
	return &Engine{
		&router{},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groups := e.router.groups
	for _, g := range groups {
		for name, methodHandle := range g.handlerMap {
			url := "/" + g.groupName + name
			if r.RequestURI == url {
				ctx := &Context{
					W: w,
					R: r,
				}
				_, ok := methodHandle["ANY"]
				if ok {
					methodHandle["ANY"](ctx)
					return
				}
				method := r.Method
				handler, ok := methodHandle[method]
				if ok {
					handler(ctx)
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintln(w, method+" not allowed ")
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, r.RequestURI+" not found ")
				return
			}
		}
	}

}

func (e *Engine) Run() {

	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}

}
