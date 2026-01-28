package storage

import (
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	rdb *redis.Client
}

func NewRedisProvider(rdb *redis.Client) *RedisProvider {
	return &RedisProvider{rdb: rdb}
}

func (r *RedisProvider) GetIP(ip string) (*models.IPstate, bool, error) {
	return RedisGet[models.IPstate](r.rdb, "ip", ip)
}

func (r *RedisProvider) SetIP(ip string, state *models.IPstate, ttl time.Duration) error {
	return RedisSet(r.rdb, "ip", ip, state, ttl)
}

func (r *RedisProvider) DeleteIP(ip string) error {
	return RedisDelete(r.rdb, "ip", ip)
}

func (r *RedisProvider) GetToken(token string) (*models.TokenState, bool, error) {
	return RedisGet[models.TokenState](r.rdb, "token", token)
}

func (r *RedisProvider) SetToken(token string, state *models.TokenState, ttl time.Duration) error {
	return RedisSet(r.rdb, "token", token, state, ttl)
}

func (r *RedisProvider) DeleteToken(token string) error {
	return RedisDelete(r.rdb, "token", token)
}
