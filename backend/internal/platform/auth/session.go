package auth

import (
	"sync"
	"time"
)

type MemoryRateLimiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	attempts map[string][]time.Time
}

func NewMemoryRateLimiter(max int, window time.Duration) *MemoryRateLimiter {
	return &MemoryRateLimiter{max: max, window: window, attempts: map[string][]time.Time{}}
}

func (l *MemoryRateLimiter) Allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)
	items := l.attempts[key]
	kept := items[:0]
	for _, t := range items {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.max {
		l.attempts[key] = kept
		return false, l.window - now.Sub(kept[0])
	}
	kept = append(kept, now)
	l.attempts[key] = kept
	return true, 0
}
