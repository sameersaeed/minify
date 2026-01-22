package limiter

import (
	"math"
	"sync"
	"time"
)

type Bucket struct {
	capacity      float64
	tokens        float64
	rate          float64
	lastAccessed  time.Time
	cooldown      time.Duration
	cooldownUntil time.Time
	mu            sync.Mutex
}

type Limiter struct {
	buckets    map[string]*Bucket
	mu         sync.Mutex
	maxBuckets int
}

type RateConfig struct {
	Rate     float64       // allowed rate in tokens/sec
	Capacity float64       // burst size (max requests sendable at once)
	Cooldown time.Duration // duration in seconds to block requests once rate limit is hit
}

var Rates = struct {
	Authenticated, Anonymous RateConfig
}{
	Authenticated: RateConfig{Rate: 0.50, Capacity: 10, Cooldown: 60 * time.Second},
	Anonymous:     RateConfig{Rate: 0.33, Capacity: 5, Cooldown: 120 * time.Second},
}

// NewBucket creates a new bucket to track user tokens
func NewBucket(cfg RateConfig) *Bucket {
	now := time.Now()
	return &Bucket{
		capacity:     cfg.Capacity,
		tokens:       cfg.Capacity,
		rate:         cfg.Rate,
		lastAccessed: now,
		cooldown:     cfg.Cooldown,
	}
}

// NewLimiter creates a new limiter to map user keys to token buckets
func NewLimiter(maxBuckets int) *Limiter {
	return &Limiter{
		buckets:    make(map[string]*Bucket),
		maxBuckets: maxBuckets,
	}
}

// Allow checks and consumes rate-limit capacity for the given requester
func (l *Limiter) Allow(key string, cfg RateConfig) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.buckets) >= l.maxBuckets {
		l.CleanupOldBuckets()
	}

	b, ok := l.buckets[key]
	if !ok { // create bucket if it doesn't exist yet
		b = NewBucket(cfg)
		l.buckets[key] = b
	}

	return b.Consume(time.Now())
}

// Allow checks whether a request is allowed for the bucket and consumes the token if so
func (b *Bucket) Consume(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// still in cool-down, reject
	if now.Before(b.cooldownUntil) {
		return false
	}

	// refill tokens
	elapsed := now.Sub(b.lastAccessed).Seconds() * b.rate
	b.tokens = math.Min(b.capacity, b.tokens+elapsed)
	b.lastAccessed = now

	// reduce for every call made
	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}

	b.cooldownUntil = now.Add(b.cooldown) // cooldown on token limit
	return false
}

func (l *Limiter) CleanupOldBuckets() {
	expiration := 30 * time.Minute

	for k, b := range l.buckets {
		// delete old buckets
		b.mu.Lock()
		last := b.lastAccessed
		b.mu.Unlock()

		if time.Since(last) >= expiration {
			delete(l.buckets, k)
		}

		// can stop once under limit
		if len(l.buckets) < l.maxBuckets {
			break
		}
	}
}
