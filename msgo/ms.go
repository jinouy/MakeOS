package msgo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type routerGroup struct {
	groupName          string
	handlerMap         map[string]map[string]HandlerFunc
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap   map[string][]string
	treeNode           *treeNode
	middlewares        []MiddlewareFunc
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(name string, method string, h HandlerFunc, ctx *Context) {
	// 组通用前置中间件
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
	// 组路由级别
	middlewareFuncs := r.middlewaresFuncMap[name][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			h = middlewareFunc(h)
		}
	}
	h(ctx)
	// 组通用后置中间件
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{
		groupName:          name,
		handlerMap:         make(map[string]map[string]HandlerFunc),
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
		handlerMethodMap:   make(map[string][]string),
		treeNode:           &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, g)
	return g
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	_, ok := r.handlerMap[name]
	if !ok {
		r.handlerMap[name] = make(map[string]HandlerFunc)
		r.middlewaresFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	_, ok = r.handlerMap[method]
	if ok {
		panic("有重复的路由")
	}
	r.handlerMap[name][method] = handlerFunc
	r.middlewaresFuncMap[name][method] = append(r.middlewaresFuncMap[name][method], middlewareFunc...)
	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, ANY, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Get(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Post(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPatch, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodDelete, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Put(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPut, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Options(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodOptions, handlerFunc, middlewareFunc...)
}
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodHead, handlerFunc, middlewareFunc...)
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

	e.HttpRequestHandle(w, r)
}

func (e *Engine) HttpRequestHandle(w http.ResponseWriter, r *http.Request) {

	for _, g := range e.routerGroups {
		routerName := SubStringLast(r.RequestURI, "/"+g.groupName)
		node := g.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			ctx := &Context{
				W: w,
				R: r,
			}
			handle, ok := g.handlerMap[node.routerName][ANY]
			if ok {
				g.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			method := r.Method
			handle, ok = g.handlerMap[node.routerName][method]
			if ok {
				g.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, method+" not allowed ")
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, r.RequestURI+" not found ")
		return
	}
}

func (e *Engine) Run() {

	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}

}
