package storage

import (
	"sync"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
)

type MemoryProvider struct {
	mu           sync.RWMutex
	ipStorage    map[string]*models.IPstate
	tokenStorage map[string]*models.TokenState
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		ipStorage:    make(map[string]*models.IPstate),
		tokenStorage: make(map[string]*models.TokenState),
	}
}

func (m *MemoryProvider) GetIP(ip string) (*models.IPstate, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.ipStorage[ip]
	if !exists {
		return nil, false, nil
	}
	return state, true, nil
}

func (m *MemoryProvider) SetIP(ip string, state *models.IPstate, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ipStorage[ip] = state
	return nil
}

func (m *MemoryProvider) DeleteIP(ip string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.ipStorage, ip)
	return nil
}

func (m *MemoryProvider) GetToken(token string) (*models.TokenState, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.tokenStorage[token]
	if !exists {
		return nil, false, nil
	}
	return state, true, nil
}

func (m *MemoryProvider) SetToken(token string, state *models.TokenState, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokenStorage[token] = state
	return nil
}

func (m *MemoryProvider) DeleteToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tokenStorage, token)
	return nil
}
