package router

import "net/http"

type Router struct {
	routes map[string]http.Handler
}

func New() *Router {
	return &Router{routes: map[string]http.Handler{}}
}

func (rt *Router) Handle(path string, h http.Handler) {
	rt.routes[path] = h
}

func (rt *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	rt.Handle(path, http.HandlerFunc(f))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := rt.routes[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}
