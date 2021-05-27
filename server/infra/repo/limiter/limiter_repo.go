package limiter

import (
	"github.com/orca-zhang/lrucache"
	"github.com/poonman/entry-task/server/domain/aggr/ratelimiter"
)

type repo struct {
	cache *lrucache.LRUCache
}

func (r *repo) Get(username string) (lim *ratelimiter.RateLimiter) {
	ret, ok := r.cache.Get(username)
	if !ok {
		return nil
	}

	lim = ret.(*ratelimiter.RateLimiter)
	return
}

func (r *repo) Save(lim *ratelimiter.RateLimiter) {
	r.cache.Put(lim.Username, lim)
}

func NewRepo() ratelimiter.Repo {
	r := &repo{
		cache: lrucache.New(10000),
	}

	return r
}
