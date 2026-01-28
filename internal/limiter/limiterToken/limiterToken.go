package limiterToken

import (
	"errors"
	"fmt"
	"os"
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
			WindowStart: time.Now(),
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

	if !getTokenState.BannedUntil.IsZero() && time.Now().After(getTokenState.BannedUntil) {
		storage.SetTokenState(token,
			&models.TokenState{
				Count:       1,
				WindowStart: time.Now(),
				BannedUntil: time.Time{},
			})
		return true, nil
	}

	if getTokenState.Count >= 5 {
		if !getTokenState.BannedUntil.IsZero() && time.Now().Before(getTokenState.BannedUntil) {
			return false, errors.New("Token blocked due to too many requests")
		}

		duration := os.Getenv("timeout_token_block_inSeconds")
		seconds, err := time.ParseDuration(duration + "s")
		if err != nil {
			fmt.Println("Error parsing duration:", err)
		}
		storage.SetTokenState(token,
			&models.TokenState{
				Count:       getTokenState.Count,
				WindowStart: getTokenState.WindowStart,
				BannedUntil: time.Now().Add(seconds),
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

func CheckTokenisValid(token string) bool {
	apiKey := os.Getenv("API_KEY")
	if token != apiKey {
		return false
	}
	return true
}
