package main

import (
	"context"
	"fmt"
	"httpServer/internal/http1"
	"httpServer/internal/metrics"
	"httpServer/internal/router"
	"httpServer/internal/tracing"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const addr = ":8080"

func main() {
	r := router.New()
	shutdown, err := tracing.Setup(context.Background())
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("hello").Start(r.Context(), "hello")
		defer span.End()
		fmt.Fprintf(w, "hello")
	})
	r.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("bye").Start(r.Context(), "bye")
		defer span.End()
		fmt.Fprintf(w, "bye")
	})
	r.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("users").Start(r.Context(), "users.list")
		defer span.End()

		fmt.Fprint(w, "list users")
	})

	r.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.Tracer("users").Start(r.Context(), "users.create")
		defer span.End()
		fmt.Fprintf(w, "create user")
	})
	r.HandleFunc("GET /metrics", metrics.Handler)
	handler := otelhttp.NewHandler(metrics.Count(r), "httpserver",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
	if err := http1.ListenAndServe(addr, handler); err != nil {
		panic(err)
	}
}
