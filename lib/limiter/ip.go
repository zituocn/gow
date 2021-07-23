package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

// IPLimiter 基于 ip的限流
type IPLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPLimiter return a new ip limiter
//	每r秒发生多少个事件
//	最大缓存b个事件
func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// addIP add ip to map
//	use sync.RwMutex
func (m *IPLimiter) addIP(ip string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	limiter := rate.NewLimiter(m.r, m.b)
	m.ips[ip] = limiter
	return limiter
}

// GetLimiter return *rate.Limiter
func (m *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	m.mu.Lock()
	limiter, ok := m.ips[ip]
	if !ok {
		m.mu.Unlock()
		return m.addIP(ip)
	}
	m.mu.Unlock()
	return limiter
}
