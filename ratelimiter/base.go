package ratelimiter

// RateLimiter as a base rate limiter interface
type RateLimiter interface {
	Allow() bool
}
