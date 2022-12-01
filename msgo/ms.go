package msgo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	g := &routerGroup{
		groupName:        name,
		handlerMap:       make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode:         &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, g)
	return g
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc) {
	_, ok := r.handlerMap[name]
	if !ok {
		r.handlerMap[name] = make(map[string]HandlerFunc)
	}
	_, ok = r.handlerMap[method]
	if ok {
		panic("有重复的路由")
	}
	r.handlerMap[name][method] = handlerFunc
	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], name)
	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc) {

	r.handle(name, ANY, handlerFunc)
}

func (r *routerGroup) Handle(name string, method string, handlerFunc HandlerFunc) {

	r.handle(name, method, handlerFunc)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodGet, handlerFunc)
}
func (r *routerGroup) Post(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPost, handlerFunc)
}
func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPatch, handlerFunc)
}
func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodDelete, handlerFunc)
}
func (r *routerGroup) Put(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPut, handlerFunc)
}
func (r *routerGroup) Options(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodOptions, handlerFunc)
}
func (r *routerGroup) Head(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodHead, handlerFunc)
}

type routerGroup struct {
	groupName        string
	handlerMap       map[string]map[string]HandlerFunc
	handlerMethodMap map[string][]string
	treeNode         *treeNode
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

	for _, g := range e.routerGroups {
		routerName := SubStringLast(r.RequestURI, "/"+g.groupName)
		node := g.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			ctx := &Context{
				W: w,
				R: r,
			}
			handle, ok := g.handlerMap[node.routerName]["ANY"]
			if ok {
				handle(ctx)
				return
			}
			method := r.Method
			handler, ok := g.handlerMap[node.routerName][method]
			if ok {
				handler(ctx)
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
