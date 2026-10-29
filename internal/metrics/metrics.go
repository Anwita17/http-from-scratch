package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var total_req atomic.Int64

func Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		total_req.Add(1)
		next.ServeHTTP(w, r)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.1")
	fmt.Println(w, "http_requests_total %d\n", total_req.Load())
	fmt.Fprintf(w, "http_requests_total %d\n", total_req.Load())
}
