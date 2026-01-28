package main

import (
	"fmt"
	"net/http"

	middlewares "github.com/gabrielpgava/rate-limiter-fullcycle/internal/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	r:= chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middlewares.RateLimiterMiddle)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	http.ListenAndServe(":8080", r)
}