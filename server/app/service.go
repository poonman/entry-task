package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/kv"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/domain/aggr/session"
)

type Service struct {
	sessionRepo session.Repo
	kvRepo    kv.Repo
	quotaRepo quota.Repo
}

func NewService(
	sessionRepo session.Repo,
	kvRepo kv.Repo,
	quotaRepo quota.Repo,
) *Service {
	return &Service{
		sessionRepo: sessionRepo,
		kvRepo:    kvRepo,
		quotaRepo: quotaRepo,
	}
}
