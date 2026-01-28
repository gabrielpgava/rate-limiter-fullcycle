package main

import (
	"net/http"
	"os"
	"strings"

	middlewares "github.com/gabrielpgava/rate-limiter-fullcycle/internal/middleware"
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func main() {

	storageProvider := strings.ToLower(os.Getenv("STORAGE_PROVIDER"))
	switch storageProvider {
	case "redis":
		addr := envOr("REDIS_ADDR", "localhost:6379")
		password := os.Getenv("REDIS_PASSWORD")
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
		storage.Use(storage.NewRedisProvider(rdb))
	default:
		storage.Use(storage.NewMemoryProvider())
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middlewares.RateLimiterMiddle)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Your request has arrived successfully!"))
	})

	http.ListenAndServe(":8080", r)
}

func envOr(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
