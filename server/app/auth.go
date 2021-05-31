package app

import (
	"context"
	"github.com/poonman/entry-task/dora/status"
	uuid "github.com/satori/go.uuid"
	"strings"
)

var (
	ErrUnauthenticated = status.New(status.Unauthenticated, "unauthenticated due to an invalid token")
)

func (s *Service) Login(ctx context.Context, username, password string) (token string, err error) {

	acc, err := s.accountRepo.Get(ctx, username)
	if err != nil {
		return
	}

	err = acc.ValidatePassword(password)
	if err != nil {
		return
	}

	// FIXME: here is a simple way to generate a token.
	token = getToken()

	// FIXME: just save in memory without ttl, the best way is saving to redis if need
	err = s.sessionRepo.Save(username, token)
	if err != nil {
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
