package ratelimiter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindowRateLimiter(t *testing.T) {
	t.Parallel()

	rateLimiter := NewFixedWindowsRateLimiter(context.Background(), 3, 100*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(150 * time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())
}

func TestFixedWindowLimiterWithGoroutines(t *testing.T) {
	t.Parallel()

	var goroutinesNumber uint32 = 10
	rateLimiter := NewFixedWindowsRateLimiter(context.Background(), 10, time.Second)

	wg := sync.WaitGroup{}
	wg.Add(int(goroutinesNumber))
	for i := range goroutinesNumber {
		println(i)
		go func() {
			defer wg.Done()
			assert.True(t, rateLimiter.Allow())
		}()
	}

	wg.Wait()
	assert.False(t, rateLimiter.Allow())
}
