package app

import (
	"github.com/poonman/entry-task/server/domain/aggr/account"
	"github.com/poonman/entry-task/server/domain/aggr/kv"
	"github.com/poonman/entry-task/server/domain/aggr/quota"
	"github.com/poonman/entry-task/server/domain/aggr/session"
	"github.com/poonman/entry-task/server/domain/factory"
)

type Service struct {
	sessionRepo session.Repo
	kvRepo      kv.Repo
	quotaRepo   quota.Repo
	accountRepo account.Repo

	factory *factory.Factory
}

func NewService(
	accountRepo account.Repo,
	sessionRepo session.Repo,
	kvRepo kv.Repo,
	quotaRepo quota.Repo,
	factory *factory.Factory,
) *Service {
	return &Service{
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		kvRepo:      kvRepo,
		quotaRepo:   quotaRepo,
		factory:     factory,
	}
}
