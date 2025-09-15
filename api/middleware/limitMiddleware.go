package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Limiter interface {
	Allow() bool
}

var (
	limiters                 = make(map[string]Limiter)
	mutex                    sync.RWMutex
	ErrLimiterNotInitialized = errors.New("limiter not initialized")
	defaultLimiter           = NewTokenBucketLimiter
)

type TokenBucketLimiter struct {
	capacity   int
	tokens     int
	rate       time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucketLimiter(capacity int, rate time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now().Add(-rate * time.Duration(capacity)),
	}
}
func (l *TokenBucketLimiter) Allow() bool {
	var err error
	if l == nil {
		err = ErrLimiterNotInitialized
		panic(err)
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastRefill)
	tokensToAdd := int(elapsed / l.rate)
	l.tokens += tokensToAdd
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	//这里如果频繁点击，超过rate，就会一直无法生成令牌
	l.lastRefill = now
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}

type SlidingWindowLimiter struct {
	windowSize time.Duration
	limit      int
	requests   []time.Time
	mu         sync.Mutex
}

func NewSlidingWindowLimiter(windowSize time.Duration, limit int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize: windowSize,
		limit:      limit,
		requests:   make([]time.Time, 0),
	}
}

func (l *SlidingWindowLimiter) Allow() bool {
	if l == nil {
		err := ErrLimiterNotInitialized
		panic(err)
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	//移除过期请求
	for len(l.requests) > 0 && now.Sub(l.requests[0]) > l.windowSize {
		l.requests = l.requests[1:]
	}
	if len(l.requests) >= l.limit {
		return false
	}
	l.requests = append(l.requests, now)
	return true
}

func LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		mutex.RLock()
		limiter, ok := limiters[clientIP]
		mutex.RUnlock()
		if !ok {
			mutex.Lock()
			limiter = defaultLimiter(5, time.Millisecond*500)
			limiters[clientIP] = limiter
			mutex.Unlock()
		}
		if !limiter.Allow() {
			c.Abort()
			c.JSON(429, gin.H{
				"msg": "Too Many Requests, Please Try Again Later",
			})
			return
		}
		c.Next()
	}
}

// 申请令牌（返回true表示成功）
func acquireToken(redisCli *redis.Client, tokenType string, count int) bool {
	key := fmt.Sprintf("global_token_bucket:%s", tokenType)
	now := time.Now().UnixMilli() // 当前时间戳（毫秒）

	// 调用Lua脚本
	res, err := redisCli.Eval(context.Background(),
		tokenBucketScript, // 上文的Lua脚本
		[]string{key},     // KEYS[1]
		count, now,        // ARGV[1]=申请数量，ARGV[2]=当前时间
	).Int64()

	return err == nil && res == 1
}

const (
	tokenBucketScript = `	
-- 输入参数：
-- KEYS[1]：令牌桶键名（如global_token_bucket:common）
-- ARGV[1]：申请的令牌数量（如1）
-- ARGV[2]：当前时间戳（毫秒，由客户端传入，避免Redis时间不一致）

-- 1. 获取当前令牌桶状态
local bucket = redis.call('HMGET', KEYS[1], 'capacity', 'rate', 'last_refresh_time', 'current_tokens')
local capacity = tonumber(bucket[1])
local rate = tonumber(bucket[2])  -- 单位：个/毫秒
local last_time = tonumber(bucket[3])
local current = tonumber(bucket[4] or 0)

-- 2. 计算时间差，动态生成新令牌
local  now = tonumber(ARGV[2])
local delta = now - last_time  -- 毫秒差值

if delta > 0 then
-- 生成新令牌 = 时间差 * 速率（最多不超过桶容量）
local new_tokens = delta * rate
current = math.min(current + new_tokens, capacity)
-- 更新上次刷新时间
last_time = now
end

-- 3. 检查是否有足够令牌分配
local request = tonumber(ARGV[1])
if current >= request then
current = current - request  -- 分配令牌
-- 更新令牌桶状态
redis.call('HMSET', KEYS[1],
'last_refresh_time', last_time,
'current_tokens', current)
return 1  -- 分配成功
else
-- 令牌不足，不更新状态
return 0  -- 分配失败
end
`
)
