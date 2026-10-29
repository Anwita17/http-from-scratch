package main

import (
	"fmt"
	"httpServer/internal/http1"
	"httpServer/internal/metrics"
	"httpServer/internal/router"
	"net/http"
)

const addr = ":8080"

func main() {
	r := router.New()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	r.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "bye")
	})
	r.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "list users")
	})
	r.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "create user")
	})
	r.HandleFunc("GET /metrics", metrics.Handler)
	if err := http1.ListenAndServe(addr, metrics.Count(r)); err != nil {
		panic(err)
	}
}
