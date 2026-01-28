//go:build integration

package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
	"github.com/redis/go-redis/v9"
)

func TestRateLimiterMiddle_Redis429Message(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis not available at %s: %v", addr, err)
	}

	if err := rdb.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush redis: %v", err)
	}

	storage.Use(storage.NewRedisProvider(rdb))
	t.Setenv("ENABLE_IP_LIMITER", "true")
	t.Setenv("ENABLE_TOKEN_LIMITER", "false")
	t.Setenv("max_request_ip_per_second", "1")
	t.Setenv("timeout_ip_block_inSeconds", "5")

	handler := RateLimiterMiddle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req1 := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req1.RemoteAddr = "127.0.0.1:1234"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("expected 200 on first request, got %d", rr1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req2.RemoteAddr = "127.0.0.1:1234"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on second request, got %d", rr2.Code)
	}

	expected := "you have reached the maximum number of requests or actions allowed within a certain time frame"
	if rr2.Body.String() != expected {
		t.Fatalf("unexpected body: %q", rr2.Body.String())
	}

}
