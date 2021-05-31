package redis

import (
	"github.com/poonman/entry-task/server/infra/config"
	"time"
)
import "github.com/gomodule/redigo/redis"

func NewRedisPool(conf *config.Config) *redis.Pool {
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisConfig.Address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", conf.RedisConfig.Password); err != nil {
				c.Close()
				return nil, err
			}
			//if _, err := c.Do("SELECT", db); err != nil {
			//  c.Close()
			//  return nil, err
			//}
			return c, nil
		},
		DialContext:     nil,
		TestOnBorrow:    nil,
		MaxIdle:         conf.RedisConfig.MaxIdle,
		MaxActive:       conf.RedisConfig.MaxActive,
		IdleTimeout:     100 * time.Second,
		Wait:            false,
		MaxConnLifetime: 500 * time.Second,
	}

	return pool
}
