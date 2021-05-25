package kv

import "fmt"

func NewKey(uid uint64, originKey string) (key string) {
	return fmt.Sprintf("%d:%s", uid, originKey)
}
