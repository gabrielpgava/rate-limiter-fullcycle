package limiterToken

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)

func CheckTokenLimit(token string) (bool, error) {

	isValid := CheckTokenisValid(token)
	if !isValid {
		return false, errors.New("IsInvalid")
	}

	getTokenState, exists := storage.GetTokenState(token)
	if !exists {
		newState := &models.TokenState{
			Count:       1,
			WindowStart: now(),
			BannedUntil: time.Time{},
		}
		storage.SetTokenState(token, newState)
		return true, nil
	}

	fmt.Printf("token %s,count: %d, window_start: %v, banned_until: %v\n",
		token,
		getTokenState.Count,
		getTokenState.WindowStart,
		getTokenState.BannedUntil)

	if !getTokenState.BannedUntil.IsZero() {
		if now().Before(getTokenState.BannedUntil) {
			return false, errors.New("Token blocked due to too many requests")
		}
		storage.SetTokenState(token,
			&models.TokenState{
				Count:       1,
				WindowStart: now(),
				BannedUntil: time.Time{},
			})
		return true, nil
	}

	if now().Sub(getTokenState.WindowStart) >= time.Second {
		storage.SetTokenState(token,
			&models.TokenState{
				Count:       1,
				WindowStart: now(),
				BannedUntil: time.Time{},
			})
		return true, nil
	}

	maxRequestsStr := os.Getenv("max_request_token_per_second")
	maxRequests := 5
	if maxRequestsStr != "" {
		if v, err := strconv.Atoi(maxRequestsStr); err == nil {
			maxRequests = v
		}
	}

	if getTokenState.Count >= maxRequests {
		duration := os.Getenv("timeout_token_block_inSeconds")
		if duration == "" {
			duration = "10"
		}
		seconds, err := time.ParseDuration(duration + "s")
		if err != nil {
			fmt.Println("Error parsing duration:", err)
		}
		storage.SetTokenState(token,
			&models.TokenState{
				Count:       getTokenState.Count,
				WindowStart: getTokenState.WindowStart,
				BannedUntil: now().Add(seconds),
			})
		return false, errors.New("Token blocked due to too many requests")

	}

	storage.SetTokenState(token,
		&models.TokenState{
			Count:       getTokenState.Count + 1,
			WindowStart: getTokenState.WindowStart,
			BannedUntil: getTokenState.BannedUntil,
		})

	return true, nil
}

var now = time.Now

func CheckTokenisValid(token string) bool {
	apiKey := os.Getenv("API_KEY")
	if token != apiKey {
		return false
	}
	return true
}
