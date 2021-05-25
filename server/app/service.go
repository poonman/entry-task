package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/kv"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
)

type Service struct {
	kvRepo    kv.Repo
	quotaRepo quota.Repo
}

func NewService(
	kvRepo kv.Repo,
	quotaRepo quota.Repo,
) *Service {
	return &Service{
		kvRepo:    kvRepo,
		quotaRepo: quotaRepo,
	}
}
