package app

import "github.com/poonman/entry-task/server/domain/aggr/quota"

func (s *Service) GetQuota(uid int64) (q *quota.Quota, err error) {
	q, err = s.quotaRepo.Get(uid)
	if err != nil {

	}

	return
}
