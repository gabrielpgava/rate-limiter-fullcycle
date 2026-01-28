package limiterIP

import (
	"testing"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)

func TestCheckIPLimit_WindowAndBan(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("max_request_ip_per_second", "2")
	t.Setenv("timeout_ip_block_inSeconds", "5")

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	now = func() time.Time { return current }
	t.Cleanup(func() { now = time.Now })

	allowed, err := CheckIPLimit("10.0.0.1")
	if err != nil || !allowed {
		t.Fatalf("expected allow on first request, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = CheckIPLimit("10.0.0.1")
	if err != nil || !allowed {
		t.Fatalf("expected allow on second request, got allowed=%v err=%v", allowed, err)
	}

	allowed, err = CheckIPLimit("10.0.0.1")
	if err == nil || allowed {
		t.Fatalf("expected block on third request, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(3 * time.Second)
	allowed, err = CheckIPLimit("10.0.0.1")
	if err == nil || allowed {
		t.Fatalf("expected still blocked during ban, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(6 * time.Second)
	allowed, err = CheckIPLimit("10.0.0.1")
	if err != nil || !allowed {
		t.Fatalf("expected allow after ban expired, got allowed=%v err=%v", allowed, err)
	}
}

func TestCheckIPLimit_WindowReset(t *testing.T) {
	storage.Use(storage.NewMemoryProvider())
	t.Setenv("max_request_ip_per_second", "1")
	t.Setenv("timeout_ip_block_inSeconds", "5")

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	current := base
	now = func() time.Time { return current }
	t.Cleanup(func() { now = time.Now })

	allowed, err := CheckIPLimit("10.0.0.2")
	if err != nil || !allowed {
		t.Fatalf("expected allow on first request, got allowed=%v err=%v", allowed, err)
	}

	current = base.Add(1100 * time.Millisecond)
	allowed, err = CheckIPLimit("10.0.0.2")
	if err != nil || !allowed {
		t.Fatalf("expected allow after window reset, got allowed=%v err=%v", allowed, err)
	}
}
