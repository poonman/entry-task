package kv

import (
	"fmt"
	"github.com/poonman/entry-task/dora/status"
)

func NewKey(uid uint64, originKey string) (key string) {
	return fmt.Sprintf("%d:%s", uid, originKey)
}

func ValidateKey(key string) (err error) {
	if len(key) == 0 {
		err = status.New(status.BadRequest, "empty key")
		return
	}

	return
}