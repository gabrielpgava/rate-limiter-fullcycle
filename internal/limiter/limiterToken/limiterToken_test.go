package limiterToken

import (
	"testing"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)

func TestCheckTokenLimit_WindowAndBan(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("API_KEY", "token_ok")
	t.Setenv("max_request_token_per_second", "2")
	t.Setenv("timeout_token_block_inSeconds", "5")

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	now = func() time.Time { return current }
	t.Cleanup(func() { now = time.Now })

	allowed, err := CheckTokenLimit("token_ok")
	if err != nil || !allowed {
		t.Fatalf("expected allow on first request, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = CheckTokenLimit("token_ok")
	if err != nil || !allowed {
		t.Fatalf("expected allow on second request, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = CheckTokenLimit("token_ok")
	if err == nil || allowed {
		t.Fatalf("expected block on third request, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(3 * time.Second)
	allowed, err = CheckTokenLimit("token_ok")
	if err == nil || allowed {
		t.Fatalf("expected still blocked during ban, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(6 * time.Second)
	allowed, err = CheckTokenLimit("token_ok")
	if err != nil || !allowed {
		t.Fatalf("expected allow after ban expired, got allowed=%v err=%v", allowed, err)
	}
}

func TestCheckTokenLimit_WindowReset(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("API_KEY", "token_ok")
	t.Setenv("max_request_token_per_second", "1")
	t.Setenv("timeout_token_block_inSeconds", "5")

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	now = func() time.Time { return current }
	t.Cleanup(func() { now = time.Now })

	allowed, err := CheckTokenLimit("token_ok")
	if err != nil || !allowed {
		t.Fatalf("expected allow on first request, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(1100 * time.Millisecond)
	allowed, err = CheckTokenLimit("token_ok")
	if err != nil || !allowed {
		t.Fatalf("expected allow after window reset, got allowed=%v err=%v", allowed, err)
	}
}
