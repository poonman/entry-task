package excp

import (
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/idl/kv"
)

var (
	ErrIncorrectUsernameOrPassword = status.New(status.Code(kv.CODE_INCORRECT_USERNAME_OR_PASSWORD), "incorrect username or password")
	ErrKeyNotExist                 = status.New(status.Code(kv.CODE_KEY_NOT_EXIST), "key not exist")
)

func IsBuzError(err error) bool {
	st, ok := err.(*status.Status)
	if !ok {
		return false
	}

	if st.Code > status.Max {
		return true
	}

	return false
}

func Error2Status(err error) *kv.Status {
	st := status.Error2Status(err)

	s := &kv.Status{
		Code:    kv.CODE(st.Code),
		Message: st.Message,
	}

	return s
}
