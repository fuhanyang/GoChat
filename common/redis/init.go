package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func Init(config *RedisConfig) *redis.Pool {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				config.Network,
				addr,
				redis.DialPassword(config.Password),
				redis.DialDatabase(config.Db))
		},
	}
	return pool
}
