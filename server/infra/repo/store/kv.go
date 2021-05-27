package store

import (
	"github.com/gomodule/redigo/redis"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/domain/aggr/kv"
)

type repo struct {
	pool *redis.Pool
}

func (r *repo) Set(key, value string) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	if err != nil {
		err = status.New(status.InternalServerError, err.Error())
		return
	}

	return
}

func (r *repo) Get(key string) (value string, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			value = ""
			err = status.ErrNotFound
			return
		}

		err = status.New(status.InternalServerError, err.Error())
		return
	}

	return
}

func NewRepo(pool *redis.Pool) kv.Repo {
	r := &repo{
		pool: pool,
	}

	return r
}