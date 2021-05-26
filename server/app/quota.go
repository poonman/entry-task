package app

import (
	"context"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
)

func (s *Service) GetQuota(ctx context.Context, username string) (q *quota.Quota, err error) {
	q, err = s.quotaRepo.Get(ctx, username)
	if err != nil {

	}

	return
}
