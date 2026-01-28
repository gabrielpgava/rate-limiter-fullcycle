package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)

func TestRateLimiterMiddle_InvalidToken(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("API_KEY", "expected")
	t.Setenv("ENABLE_TOKEN_LIMITER", "true")

	handler := RateLimiterMiddle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req.Header.Set("API_KEY", "wrong")
	req.RemoteAddr = "127.0.0.1:1234"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestRateLimiterMiddle_TokenOverridesIP(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("API_KEY", "token_ok")
	t.Setenv("ENABLE_IP_LIMITER", "true")
	t.Setenv("ENABLE_TOKEN_LIMITER", "true")
	t.Setenv("max_request_ip_per_second", "1")
	t.Setenv("max_request_token_per_second", "2")
	t.Setenv("timeout_ip_block_inSeconds", "5")
	t.Setenv("timeout_token_block_inSeconds", "5")

	handler := RateLimiterMiddle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		req.Header.Set("API_KEY", "token_ok")
		req.RemoteAddr = "127.0.0.1:1234"
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 on request %d, got %d", i+1, rr.Code)
		}
	}
}
