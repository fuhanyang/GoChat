package authentication

import (
	"common/bloomFilter"
	"common/chain"
	redis2 "common/redis"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"log"
	Redis2 "user/DAO/Redis"
	"user/Logic/UserCreate"
	"user/Models"
)

var (
	UserNameBloomFilter bloomFilter.BloomFilter
	ErrUserNameExists   = errors.New("username already exists")
)

func NewBloomFilter(bitmapLen int64, hashCount int32, db *gorm.DB) {
	log.Println("bloomFilter init start")
	UserNameBloomFilter = bloomFilter.NewLocalBloomFilter(bitmapLen, hashCount)
	Models.InitBloomFilter(UserNameBloomFilter, db)
	log.Println("bloomFilter init success")
}

// UserSignUp 注册用户
func UserSignUp(username string, password string, ip string, db *gorm.DB) (string, error) {
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), bloomFilterHandler(), registerMysqlHandler(), registerRedisHandler())
	ctx.Set("username", username)
	ctx.Set("password", password)
	ctx.Set("ip", ip)
	ctx.Set("db", db)
	ctx.Next()
	if ctx.GetError() != nil {
		return "", ctx.GetError()
	}
	accountNum, err := chain.GetToType[string](ctx, "accountNum")
	return accountNum, err
}

func bloomFilterHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		username, err := chain.GetToType[string](ctx, "username")
		if err != nil {
			ctx.Abort()
			return err
		}
		if UserNameBloomFilter.Exists(username) {
			return ErrUserNameExists
		}

		ctx.Next()

		if ctx.GetError() == nil {
			// 写入bloomFilter
			UserNameBloomFilter.Set(username)
			fmt.Println("user Sign Up Success")
		}
		return nil
	}
}

func registerMysqlHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		username, err := chain.GetToType[string](ctx, "username")
		if err != nil {
			return err
		}
		password, err := chain.GetToType[string](ctx, "password")
		if err != nil {
			return err
		}
		ip, err := chain.GetToType[string](ctx, "ip")
		if err != nil {
			return err
		}
		db, err := chain.GetToType[*gorm.DB](ctx, "db")
		if err != nil {
			return err
		}

		// 创建用户
		user, err := UserCreate.CreateUser(username, ip, password, UserCreate.CreateAccountNum())
		defer Models.ReleaseUser(user)
		if err != nil {
			return err
		}

		Models.WriteUser(db, user)
		ctx.Set("user", user)
		ctx.Set("accountNum", user.GetAccountNum())
		//ctx.Next()
		return nil
	}
}

func registerRedisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		user, err := chain.GetToType[*Models.User](ctx, "user")
		if err != nil {
			return err
		}
		args := redis.Args{user.GetAccountNum()}.AddFlat(user)
		_, err = Redis2.RedisDo(redis2.HMSET, args...)
		return err
	}
}
