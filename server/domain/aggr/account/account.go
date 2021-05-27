package account

import (
	"github.com/poonman/entry-task/server/domain/excp"
	"strings"
)

type Account struct {
	Id       int
	Username string
	Password string
}

func (a *Account) ValidatePassword(password string) (err error) {
	if strings.Compare(a.Password, password) == 0 {
		return nil
	}

	err = excp.ErrIncorrectUsernameOrPassword

	return
}
