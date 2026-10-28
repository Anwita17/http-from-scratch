package router

import (
	"net/http"
	"sort"
	"strings"
)

type route struct {
	method  string
	path    string
	handler http.Handler
}

type Router struct {
	routes []route
}

func New() *Router {
	return &Router{}
}

func (rt *Router) Handle(pattern string, h http.Handler) {
	method, path, found := strings.Cut(pattern, " ")
	if !found {
		method, path = "", pattern
	}
	rt.routes = append(rt.routes, route{method: method, path: path, handler: h})
}

func (rt *Router) HandleFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	rt.Handle(pattern, http.HandlerFunc(f))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allowed []string
	for _, rte := range rt.routes {
		if rte.path != r.URL.Path {
			continue
		}
		if rte.method == "" || rte.method == r.Method {
			rte.handler.ServeHTTP(w, r)
			return
		}
		allowed = append(allowed, rte.method)
	}

	if len(allowed) == 0 {
		http.NotFound(w, r)
		return
	}
	sort.Strings(allowed)
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	http.Error(w, "method not llowed", http.StatusMethodNotAllowed)
}
