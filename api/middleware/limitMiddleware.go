package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
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
