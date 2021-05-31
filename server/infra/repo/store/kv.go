package store

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/domain/aggr/kv"
	"github.com/poonman/entry-task/server/infra/config"
)

type repo struct {
	pool  *redis.Pool
	conf  *config.Config
	value string
}

func (r *repo) Set(key, value string) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	if err != nil {
		err = status.New(status.InternalServerError, "set key error. "+err.Error())
		return
	}

	return
}

func (r *repo) Get(key string) (value string, err error) {

	if !r.conf.StoreRepoConfig.UseRedis {
		value = r.value
		return
	}

	conn := r.pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			value = ""
			err = status.ErrNotFound
			return
		}

		err = status.New(status.InternalServerError, "get key error. "+err.Error())
		return
	}

	return
}

func newValue(id int) (key string) {
	key = fmt.Sprintf("%dyyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy", id)
	for i := 0; i < 9; i++ {
		key += "0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy0yyyyyyyyy"
	}

	return
}

func NewRepo(conf *config.Config, pool *redis.Pool) kv.Repo {
	r := &repo{
		pool:  pool,
		conf:  conf,
		value: newValue(1),
	}

	return r
}
