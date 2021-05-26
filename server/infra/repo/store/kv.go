package store

import (
	"github.com/gomodule/redigo/redis"
	"github.com/poonman/entry-task/server/domain/aggr/kv"
)

type repo struct {
	pool *redis.Pool
}

func (r *repo) Set(key, value string) (err error) {
	panic("implement me")
}

func (r *repo) Get(key string) (value string, err error) {
	panic("implement me")
}

func NewRepo(pool *redis.Pool) kv.Repo {
	r := &repo{
		pool: pool,
	}

	return r
}