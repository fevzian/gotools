package ratelimiter

import (
	"context"
	"sync/atomic"
	"time"
)

type BucketTokenRateLimiter struct {
	capacity       uint32
	tokens         uint32
	interval       time.Duration
	lastRefillTime int64
	rate           float32
}

func NewBucketTokenRateLimiter(ctx context.Context, capacity uint32, interval time.Duration) *BucketTokenRateLimiter {
	return &BucketTokenRateLimiter{
		capacity:       capacity,
		tokens:         capacity,
		interval:       interval,
		lastRefillTime: time.Now().UnixMilli(),
		rate:           float32(interval.Milliseconds()) / float32(capacity), //millis to pass to earn a token
	}
}

func (limiter *BucketTokenRateLimiter) Allow() bool {
	now := time.Now().UnixMilli()
	elapsed := now - limiter.lastRefillTime
	// if a bucket was not used for a long time
	if elapsed >= limiter.interval.Milliseconds() {
		atomic.StoreUint32(&limiter.tokens, limiter.capacity)
		atomic.StoreInt64(&limiter.lastRefillTime, now)
	} else {
		moreTokens := float32(elapsed) / limiter.rate
		atomic.AddUint32(&limiter.tokens, uint32(moreTokens))
		atomic.StoreInt64(&limiter.lastRefillTime, now)
	}

	tokens := atomic.LoadUint32(&limiter.tokens)
	if tokens > 0 {
		atomic.AddUint32(&limiter.tokens, ^uint32(0))
		return true
	}
	return false
}
