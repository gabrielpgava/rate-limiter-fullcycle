package limiterIP

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)

func CheckIPLimit(ip string) (bool, error) {

	getIPState, exists := storage.GetIPState(ip)
	if !exists {
		newState := &models.IPstate{
			Count:       1,
			WindowStart: now(),
			BannedUntil: time.Time{},
		}
		storage.SetIPState(ip, newState)
		return true, nil
	}

	fmt.Printf("ip %s,count: %d, window_start: %v, banned_until: %v\n",
		ip,
		getIPState.Count,
		getIPState.WindowStart,
		getIPState.BannedUntil)

	if !getIPState.BannedUntil.IsZero() {
		if now().Before(getIPState.BannedUntil) {
			return false, errors.New("IP blocked due to too many requests")
		}
		storage.SetIPState(ip,
			&models.IPstate{
				Count:       1,
				WindowStart: now(),
				BannedUntil: time.Time{},
			})
		return true, nil
	}

	if now().Sub(getIPState.WindowStart) >= time.Second {
		storage.SetIPState(ip,
			&models.IPstate{
				Count:       1,
				WindowStart: now(),
				BannedUntil: time.Time{},
			})
		return true, nil
	}

	maxRequestsStr := os.Getenv("max_request_ip_per_second")
	maxRequests := 5
	if maxRequestsStr != "" {
		if v, err := strconv.Atoi(maxRequestsStr); err == nil {
			maxRequests = v
		}
	}

	if getIPState.Count >= maxRequests {
		duration := os.Getenv("timeout_ip_block_inSeconds")
		if duration == "" {
			duration = "10"
		}
		seconds, err := time.ParseDuration(duration + "s")
		if err != nil {
			fmt.Println("Error parsing duration:", err)
		}
		storage.SetIPState(ip,
			&models.IPstate{
				Count:       getIPState.Count,
				WindowStart: getIPState.WindowStart,
				BannedUntil: now().Add(seconds),
			})
		return false, errors.New("IP blocked due to too many requests")
	}

	storage.SetIPState(ip,
		&models.IPstate{
			Count:       getIPState.Count + 1,
			WindowStart: getIPState.WindowStart,
			BannedUntil: getIPState.BannedUntil,
		})

	return true, nil
}

var now = time.Now
