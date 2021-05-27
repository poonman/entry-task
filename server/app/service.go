package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/account"
	"github.com/poonman/entry-task/server/domain/aggr/kv"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/domain/aggr/session"
)

type Service struct {
	sessionRepo session.Repo
	kvRepo    kv.Repo
	quotaRepo quota.Repo
	accountRepo account.Repo
}

func NewService(
	accountRepo account.Repo,
	sessionRepo session.Repo,
	kvRepo kv.Repo,
	quotaRepo quota.Repo,
) *Service {
	return &Service{
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		kvRepo:    kvRepo,
		quotaRepo: quotaRepo,
	}
}
