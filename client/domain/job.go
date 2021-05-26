package domain

import (
	"github.com/poonman/entry-task/client/domain/aggr/user"
	"github.com/poonman/entry-task/dora/log"
)

func (s *Service) RunJob() {
	u := &user.User{
		Name:     "1",
		Password: "1",
		Token:    "",
	}

	key := "xxx"
	value := "yyy"
	err := s.kvGateway.Set(u, key, value)
	if err != nil {
		log.Errorf("Failed to set key. err:[%v]", err)
	}

	v, err := s.kvGateway.Get(u, key)
	if err != nil {
		log.Errorf("Failed to get key. err:[%v]", err)
		return
	}

	log.Debugf("v:[%v]", v)
}
