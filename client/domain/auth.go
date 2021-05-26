package domain

import "github.com/poonman/entry-task/client/domain/aggr/user"

func (s *Service) Login(u *user.User) (err error) {
	err = s.kvGateway.Login(u)
	if err != nil {
		return err
	}

	return
}

func (s *Service) SetKV(u *user.User, k, v string) (err error) {
	return s.kvGateway.Set(u, k, v)
}

func (s *Service) GetKv(u *user.User, k string)(v string, err error) {
	return s.kvGateway.Get(u, k)
}