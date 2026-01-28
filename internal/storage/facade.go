package storage

import "github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"

var activeProvider Provider
var fallbackProvider Provider

func Use(p Provider) {
	activeProvider = p
}

func provider() Provider {
	if activeProvider != nil {
		return activeProvider
	}
	if fallbackProvider == nil {
		fallbackProvider = NewMemoryProvider()
	}
	return fallbackProvider
}

func GetIPState(ip string) (*models.IPstate, bool) {
	state, exists, err := provider().GetIP(ip)
	if err != nil {
		return nil, false
	}
	return state, exists
}

func SetIPState(ip string, state *models.IPstate) {
	_ = provider().SetIP(ip, state, 0)
}

func DeleteIPState(ip string) {
	_ = provider().DeleteIP(ip)
}

func GetTokenState(token string) (*models.TokenState, bool) {
	state, exists, err := provider().GetToken(token)
	if err != nil {
		return nil, false
	}
	return state, exists
}

func SetTokenState(token string, state *models.TokenState) {
	_ = provider().SetToken(token, state, 0)
}

func DeleteTokenState(token string) {
	_ = provider().DeleteToken(token)
}
