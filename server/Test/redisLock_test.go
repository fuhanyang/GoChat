package Test

import (
	"fmt"
	redislock "github.com/jefferyjob/go-redislock"
	"time"

	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedisLock(t *testing.T) {
	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	for i := 0; i < 10; i++ {
		go func() {
			err := redislock.New(context.Background(), redisClient, "testLock").Lock()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("lock acquired")
		}()
	}

	time.Sleep(10 * time.Second)
}
