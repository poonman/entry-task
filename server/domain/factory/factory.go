package factory

import (
	"context"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/domain/aggr/ratelimiter"
)

type Factory struct {
	limiterRepo ratelimiter.Repo
	quotaRepo   quota.Repo
}

func NewFactory(
	limiterRepo ratelimiter.Repo,
	quotaRepo quota.Repo,
) *Factory {
	return &Factory{
		limiterRepo: limiterRepo,
		quotaRepo:   quotaRepo,
	}
}

func (f *Factory) TryGetRateLimiter(ctx context.Context, username string) (lim *ratelimiter.RateLimiter) {
	lim = f.limiterRepo.Get(username)
	if lim != nil {
		return lim
	}

	q, err := f.quotaRepo.Get(ctx, username)
	if err != nil {
		return nil
	}

	log.Debugf("quota:[%+v]", q)

	lim = ratelimiter.NewRateLimiter(username, q)

	f.limiterRepo.Save(lim)

	return lim
}
