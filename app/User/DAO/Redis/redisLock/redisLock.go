package redisLock

import (
	"User/DAO/Redis"
	"context"
	redislock "github.com/jefferyjob/go-redislock"
)

// NewRedisLock 创建分布式redis锁
func NewRedisLock(key string, ctx context.Context) redislock.RedisLockInter {
	lock := redislock.New(ctx, Redis.Client, key)
	if lock == nil {
		panic("redis lock get failed")
	}
	return lock
}
