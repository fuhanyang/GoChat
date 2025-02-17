package Test

import (
	Redis2 "User/DAO/Redis"
	settings "User/Settings"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestRedis(t *testing.T) {
	err := settings.Init()
	if err != nil {
		panic(err)
	}
	Redis2.Init(settings.Config.RedisConfig)
	_, err = Redis2.RedisDo(Redis2.HMSET, "test", "name", "zhangsan", "age", 20)
	if err != nil {
		panic(err)
	}
	reply, err := Redis2.RedisDo(Redis2.HMGET, "test", "name")
	if err != nil {
		panic(err)
	}
	value, _ := redis.Values(reply, err)
	s, _ := redis.String(value, nil)
	fmt.Println(s)
}
