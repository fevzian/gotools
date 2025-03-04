package ratelimiter

import (
	"context"
	"sync/atomic"
	"time"
)

type FixedWindowsRateLimiter struct {
	count           uint32
	limit           uint32
	interval        time.Duration
	windowStartTime int64
}

func NewFixedWindowsRateLimiter(ctx context.Context, limit uint32, interval time.Duration) *FixedWindowsRateLimiter {
	limiter := &FixedWindowsRateLimiter{
		count:           0,
		limit:           limit,
		interval:        interval,
		windowStartTime: time.Now().UnixMilli(),
	}
	return limiter
}

func (rateLimiter *FixedWindowsRateLimiter) Allow() bool {
	currentTimeMs := time.Now().UnixMilli()
	if (currentTimeMs - atomic.LoadInt64(&rateLimiter.windowStartTime)) >= rateLimiter.interval.Milliseconds() {
		atomic.StoreUint32(&rateLimiter.count, 0)
		atomic.StoreInt64(&rateLimiter.windowStartTime, time.Now().UnixMilli())
	}

	count := atomic.LoadUint32(&rateLimiter.count)
	if count > rateLimiter.limit {
		return false
	}

	for !atomic.CompareAndSwapUint32(&rateLimiter.count, count, count+1) {
		count = atomic.LoadUint32(&rateLimiter.count)
	}
	return count < rateLimiter.limit
}
