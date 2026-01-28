package storage

import (
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
)

var ipStorage = make(map[string]*models.IPstate)

func GetIPState(ip string) (*models.IPstate, bool) {
	state, exists := ipStorage[ip]
	return state, exists
}

func SetIPState(ip string, state *models.IPstate) {
	ipStorage[ip] = state
}

func DeleteIPState(ip string) {
	delete(ipStorage, ip)
}

var tokenStorage = make(map[string]*models.TokenState)

func GetTokenState(token string) (*models.TokenState, bool) {
	state, exists := tokenStorage[token]
	return state, exists
}

func SetTokenState(token string, state *models.TokenState) {
	tokenStorage[token] = state
}

func DeleteTokenState(token string) {
	delete(tokenStorage, token)
}
