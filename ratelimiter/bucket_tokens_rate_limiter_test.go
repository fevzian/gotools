package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketLimiter(t *testing.T) {
	t.Parallel()

	rateLimiter := NewBucketTokenRateLimiter(3, time.Second)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(350 * time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())
}

func TestTokenBucketLimiterWithGoroutines(t *testing.T) {
	t.Parallel()

	goroutinesNumber := 10
	rateLimiter := NewBucketTokenRateLimiter(uint32(goroutinesNumber), time.Second)

	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)
	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()
			assert.True(t, rateLimiter.Allow())
		}()
	}

	wg.Wait()
	assert.False(t, rateLimiter.Allow())
}

func TestTokenBucketLimiterWithCancel(t *testing.T) {
	t.Parallel()

	rateLimiter := NewBucketTokenRateLimiter(3, 500*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(550 * time.Millisecond)
	assert.True(t, rateLimiter.Allow())
}
