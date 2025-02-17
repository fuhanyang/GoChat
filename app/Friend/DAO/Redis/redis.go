package Redis

import (
	settings "Friend/Settings"
	"fmt"
	"github.com/gomodule/redigo/redis"
	redis2 "github.com/redis/go-redis/v9"
	"time"
)

var pool *redis.Pool
var Client *redis2.Client

func Init(config *settings.RedisConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	pool = &redis.Pool{
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
	Client = redis2.NewClient(&redis2.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.Db,
	})
}

func RedisDo(commandName string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()
	reply, err := conn.Do(commandName, args...)
	return reply, err
}
