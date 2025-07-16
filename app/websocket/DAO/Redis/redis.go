package Redis

import (
	"github.com/gomodule/redigo/redis"
	redis2 "github.com/redis/go-redis/v9"
)

var Pool *redis.Pool
var Client *redis2.Client

func RedisDo(commandName string, args ...interface{}) (interface{}, error) {
	conn := Pool.Get()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()
	reply, err := conn.Do(commandName, args...)
	return reply, err
}
func RedisPoolGet() redis.Conn {
	return Pool.Get()
}
func RedisPoolPut(conn redis.Conn) {
	err := conn.Close()
	if err != nil {
		panic(err)
	}
}
