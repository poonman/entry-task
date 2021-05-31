package domain

import (
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/gateway"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/client/infra/helper"
)

type Service struct {
	kvGateway gateway.KvGateway
	conf      *config.Config
}

func NewService(
	kvGateway gateway.KvGateway,
	conf *config.Config,
) *Service {

	s := &Service{
		kvGateway: kvGateway,
		conf:      conf,
	}

	return s
}

func (s *Service) Login(u *user.User) (err error) {
	err = s.kvGateway.Login(u)
	if err != nil {
		return err
	}

	return
}

func (s *Service) SetKV(u *user.User, k, v byte) (err error) {

	key := helper.NewString(k)
	value := helper.NewString(v)

	return s.kvGateway.Set(u, key, value)
}

func (s *Service) GetKV(u *user.User, k byte) (v string, err error) {

	key := helper.NewString(k)

	return s.kvGateway.Get(u, key)
}
