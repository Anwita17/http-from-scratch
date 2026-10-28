package main

import (
	"fmt"
	"httpServer/internal/http1"
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
	if err := http1.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
