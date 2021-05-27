package app

import (
	"github.com/poonman/entry-task/client/domain"
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/client/domain/aggr/userm"
	"github.com/poonman/entry-task/dora/log"
)

type Service struct {
	domainSvc *domain.Service

	userManager *userm.UserManager
}

func NewService(domainSvc *domain.Service) *Service {
	return &Service{
		domainSvc:   domainSvc,
		userManager: userm.NewUserManager(),
	}
}

func (s *Service) Login(username string) (err error) {
	u := &user.User{
		Name:     username,
		Password: username,
		Token:    "",
	}

	err = s.domainSvc.Login(u)
	if err != nil {
		return
	}

	s.userManager.AddUser(u)

	return
}

func (s *Service) WriteSecureMessage(username string) (err error) {
	u := s.userManager.GetUser(username)

	if u == nil {
		log.Warnf("no such user. username:[%s]", username)
		return
	}

	err = s.domainSvc.SetKV(u, "", "")
	if err != nil {
		log.Errorf("failed to set kv. err:[%v]", err)
		return
	}

	return
}

func (s *Service) ReadSecureMessage(username string) (err error) {
	u := s.userManager.GetUser(username)

	if u == nil {
		log.Warnf("no such user. username:[%s]", username)
		return
	}

	//var v string
	_, err = s.domainSvc.GetKv(u, "")
	if err != nil {
		log.Errorf("failed to set kv. err:[%v]", err)
		return
	}

	//log.Debugf("GetKV success. username:[%s], value:[%s]", username, v)

	return
}

func (s *Service) BenchmarkRead() {
	s.domainSvc.BenchmarkRead()
}
