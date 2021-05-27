package app

import (
	"context"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/status"
	uuid "github.com/satori/go.uuid"
	"strings"
)

var (
	ErrUnauthenticated = status.New(status.Unauthenticated, "unauthenticated")
)

func (s *Service) Login(ctx context.Context, username, password string) (token string, err error) {

	acc, err := s.accountRepo.Get(ctx, username)
	if err != nil {
		log.Errorf("failed to get account. username:[%s], err:[%v]", username, err)
		return
	}

	err = acc.ValidatePassword(password)
	if err != nil {
		return
	}

	// todo: validate username and password
	token = getToken()

	// 写入session中
	err = s.sessionRepo.Save(username, token)
	if err != nil {
		log.Errorf("Failed to save user token. err:[%v]", err)
		return
	}

	return
}

func (s *Service) Authenticate(username, token string) (err error) {
	var (
		origToken string
	)
	origToken, err = s.sessionRepo.Get(username)
	if err != nil {
		log.Errorf("Failed to get user token. err:[%v]",err)
		return
	}

	if strings.Compare(token, origToken) != 0 {
		err = ErrUnauthenticated
		return
	}

	return nil
}

func getToken() (token string) {
	var (
		u4 uuid.UUID
	)

	u4 = uuid.NewV4()

	return u4.String()
}

