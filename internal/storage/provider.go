package storage

import (
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
)

type Provider interface {
	GetIP(ip string) (*models.IPstate, bool, error)
	SetIP(ip string, state *models.IPstate, ttl time.Duration) error
	DeleteIP(ip string) error
	GetToken(token string) (*models.TokenState, bool, error)
	SetToken(token string, state *models.TokenState, ttl time.Duration) error
	DeleteToken(token string) error
}
