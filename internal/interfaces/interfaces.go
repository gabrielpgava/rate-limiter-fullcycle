package interfaces

import (
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
)

type IPStorage interface {
	Get(ip string) (*models.IPstate, bool, error)
	Set(ip string, state *models.IPstate, ttl time.Duration) error
	Delete(ip string) error
}

type TokenStorage interface {
	Get(token string) (*models.TokenState, bool, error)
	Set(token string, state *models.TokenState, ttl time.Duration) error
	Delete(token string) error
}
