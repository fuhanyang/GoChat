package redis

import (
	"context"
	"fmt"
	redislock "github.com/jefferyjob/go-redislock"
	redis2 "github.com/redis/go-redis/v9"
)

func InitRedisLock(config *RedisConfig) *redis2.Client {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := redis2.NewClient(&redis2.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.Db,
	})
	return client
}

// NewRedisLock 创建分布式redis锁
func NewRedisLock(key string, ctx context.Context, Client *redis2.Client) redislock.RedisLockInter {
	lock := redislock.New(ctx, Client, key)
	if lock == nil {
		panic("redis lock get failed")
	}
	return lock
}
