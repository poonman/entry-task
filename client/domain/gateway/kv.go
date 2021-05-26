package gateway

import "github.com/poonman/entry-task/client/domain/aggr/user"

type KvGateway interface {
	Auth(u *user.User) (err error)
	Set(u *user.User, key, value string) (err error)
	Get(u *user.User, key string) (value string, err error)
}
