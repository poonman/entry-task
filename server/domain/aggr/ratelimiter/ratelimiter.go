package ratelimiter

import (
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Username     string
	readLimiter  *rate.Limiter
	writeLimiter *rate.Limiter
}

func NewRateLimiter(username string, quota *quota.Quota) *RateLimiter {
	return &RateLimiter{
		Username:     username,
		readLimiter:  rate.NewLimiter(rate.Limit(quota.ReadQuota), max(quota.ReadQuota/10, 1)),
		writeLimiter: rate.NewLimiter(rate.Limit(quota.WriteQuota), max(quota.WriteQuota/10, 1)),
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func (r *RateLimiter) AllowRead() bool {
	return r.readLimiter.Allow()
}

func (r *RateLimiter) AllowWrite() bool {
	return r.writeLimiter.Allow()
}
