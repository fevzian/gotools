package ratelimiter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketLimiter(t *testing.T) {
	t.Parallel()

	rateLimiter := NewBucketTokenRateLimiter(context.Background(), 3, time.Second)

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
	rateLimiter := NewBucketTokenRateLimiter(context.Background(), uint32(goroutinesNumber), time.Second)

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

	ctx, cancel := context.WithCancel(context.Background())
	rateLimiter := NewBucketTokenRateLimiter(ctx, 3, 500*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	cancel()

	time.Sleep(550 * time.Millisecond)
	assert.True(t, rateLimiter.Allow())
}
