package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

// UserLimiter 基于uid的用户限流
type UserLimiter struct {
	users map[int64]*rate.Limiter
	mu    *sync.RWMutex
	r     rate.Limit
	b     int
}

// NewUserLimiter return a new user limiter
//	每r秒b个
func NewUserLimiter(r rate.Limit, b int) *UserLimiter {
	return &UserLimiter{
		users: make(map[int64]*rate.Limiter),
		mu:    &sync.RWMutex{},
		r:     r,
		b:     b,
	}
}

// addUser add uid to map
func (m *UserLimiter) addUser(uid int64) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	limiter := rate.NewLimiter(m.r, m.b)
	m.users[uid] = limiter
	return limiter
}

// GetLimiter return *rate.Limiter
func (m *UserLimiter) GetLimiter(uid int64) *rate.Limiter {
	m.mu.Lock()
	limiter, ok := m.users[uid]
	if !ok {
		m.mu.Unlock()
		return m.addUser(uid)
	}
	m.mu.Unlock()
	return limiter
}
