package quota

import (
	"sync"
	"time"
)

type Decision struct {
	Allowed    bool          `json:"allowed"`
	Remaining  int           `json:"remaining"`
	RetryAfter time.Duration `json:"retryAfter"`
}
type bucket struct {
	tokens  float64
	updated time.Time
}
type Limiter struct {
	mu               sync.Mutex
	capacity, refill float64
	buckets          map[string]bucket
	now              func() time.Time
}

func New(capacity int, refillPerSecond float64) *Limiter {
	return &Limiter{capacity: float64(capacity), refill: refillPerSecond, buckets: map[string]bucket{}, now: time.Now}
}
func (l *Limiter) Allow(key string, cost int) Decision {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b, ok := l.buckets[key]
	if !ok {
		b = bucket{tokens: l.capacity, updated: now}
	}
	b.tokens = min(l.capacity, b.tokens+now.Sub(b.updated).Seconds()*l.refill)
	b.updated = now
	if b.tokens < float64(cost) {
		l.buckets[key] = b
		missing := float64(cost) - b.tokens
		return Decision{Allowed: false, Remaining: int(b.tokens), RetryAfter: time.Duration(missing / l.refill * float64(time.Second))}
	}
	b.tokens -= float64(cost)
	l.buckets[key] = b
	return Decision{Allowed: true, Remaining: int(b.tokens)}
}
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
