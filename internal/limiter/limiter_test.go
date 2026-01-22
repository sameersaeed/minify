package limiter

import (
	"testing"
	"time"
)

func TestLimiterTokenConsumption(t *testing.T) {
	l := NewLimiter(100)
	key := "user1"
	cfg := Rates.Authenticated

	t.Log("Consume all tokens in the bucket")
	for i := 0; i < int(cfg.Capacity); i++ {
		if !l.Allow(key, cfg) {
			t.Fatalf("Expected request %d to be allowed", i+1)
		}
	}

	t.Log("Next request should fail because bucket is empty")
	if l.Allow(key, cfg) {
		t.Fatal("Expected request to be blocked after capacity reached")
	}
}

func TestLimiterTokenRefill(t *testing.T) {
	l := NewLimiter(100)
	key := "user2"
	cfg := Rates.Anonymous

	t.Log("Consume all tokens")
	for i := 0; i < int(cfg.Capacity); i++ {
		l.Allow(key, cfg)
	}

	t.Log("Next request should fail due to empty bucket")
	if l.Allow(key, cfg) {
		t.Fatal("Expected request to be blocked after capacity reached")
	}

	t.Log("Simulate token refill and cooldown expired")
	b := l.buckets[key]
	b.mu.Lock()
	b.lastAccessed = b.lastAccessed.Add(-1 * time.Second)
	b.cooldownUntil = time.Now().Add(-1 * time.Second)
	b.mu.Unlock()

	t.Log("Consume refilled tokens")
	for i := 0; i < int(cfg.Rate); i++ {
		if !l.Allow(key, cfg) {
			t.Fatalf("Expected token %d to be available after refill", i+1)
		}
	}

	t.Log("Validate that making another request after token refills gets blocked")
	if l.Allow(key, cfg) {
		t.Fatal("Expected further request to be blocked until more tokens refill")
	}
}

func TestLimiterSeparateBuckets(t *testing.T) {
	l := NewLimiter(100)
	cfg := Rates.Authenticated

	t.Log("Check separate buckets for different keys")
	if !l.Allow("user1", cfg) {
		t.Fatal("User1 should be allowed")
	}
	if !l.Allow("user2", cfg) {
		t.Fatal("User2 should be allowed")
	}
	if len(l.buckets) != 2 {
		t.Fatalf("Expected 2 buckets for 2 users, got %d", len(l.buckets))
	}
}

func TestLimiterCooldown(t *testing.T) {
	l := NewLimiter(100)
	key := "user3"
	cfg := Rates.Authenticated

	t.Log("Consume all tokens to trigger cooldown")
	for i := 0; i < int(cfg.Capacity); i++ {
		l.Allow(key, cfg)
	}

	t.Log("Next request should be blocked due to cooldown")
	if l.Allow(key, cfg) {
		t.Fatal("Expected request to be blocked during cooldown")
	}

	t.Log("Simulate cooldown expiration")
	b := l.buckets[key]
	b.mu.Lock()
	b.tokens = b.capacity
	b.cooldownUntil = time.Now().Add(-1 * time.Second) // skip cooldown time
	b.mu.Unlock()

	t.Log("Cooldown finished, request should be allowed")
	if !l.Allow(key, cfg) {
		t.Fatal("Expected request to be allowed after cooldown")
	}
}

func TestLimiterCleanupOldBuckets(t *testing.T) {
	l := NewLimiter(2)
	cfg := Rates.Anonymous

	t.Log("Create two buckets to reach max capacity")
	l.Allow("user1", cfg)
	l.Allow("user2", cfg)

	t.Log("Manually set last accessed time in the past to simulate expiration")
	for _, b := range l.buckets {
		b.mu.Lock()
		b.lastAccessed = time.Now().Add(-31 * time.Minute) // expired
		b.mu.Unlock()
	}

	t.Log("Add new bucket to trigger cleanup of old buckets")
	l.Allow("user3", cfg)

	if len(l.buckets) > 2 {
		t.Fatalf("Expected max bucket limit to be enforced after cleanup, got %d", len(l.buckets))
	}

	t.Log("Verify that new bucket exists and cleanup worked correctly")
	if _, ok := l.buckets["user3"]; !ok {
		t.Fatal("Expected new bucket 'user3' to exist after cleanup")
	}
}
