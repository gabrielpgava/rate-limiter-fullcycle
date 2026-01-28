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