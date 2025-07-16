package SearchUser

import (
	"common/chain"
	redis2 "common/redis"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"user/DAO/Redis"
	"user/Models"
)

// SearchUserByAccountNum 在数据库中通过账号查找用户信息
func SearchUserByAccountNum(accountNum string, db *gorm.DB) (*Models.User, error) {
	user := Models.NewUser()
	defer Models.ReleaseUser(user)
	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), redisHandler(), mysqlHandler())
	ctx.Set("user", user)
	ctx.Set("accountNum", accountNum)
	ctx.Set("db", db)
	err := ctx.Apply()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func redisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		// 判断在redis中是否存在该用户信息
		//先从redis中查找用户信息，如果没有再从mysql中查找
		user, err := chain.GetToType[*Models.User](ctx, "user")
		if err != nil {
			return err
		}
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}

		var values []interface{}
		reply, err := Redis.RedisDo(redis2.HGETALL, accountNum)
		if err != nil && !errors.Is(err, redis.ErrNil) {
			//产生了其他错误则要返回
			return err
		}
		values, err = redis.Values(reply, err)
		//redis中不存在该用户信息则从mysql中查找
		if errors.Is(err, redis.ErrNil) || len(values) == 0 {
			return nil
		}
		err = redis.ScanStruct(values, user)
		if err != nil {
			return err
		}
		_, err = Redis.RedisDo(redis2.EXPIRE, accountNum, 3600)
		return err
	}
}
func mysqlHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		user, err := chain.GetToType[*Models.User](ctx, "user")
		if err != nil {
			return err
		}
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}
		db, err := chain.GetToType[*gorm.DB](ctx, "db")
		if err != nil {
			return err
		}

		//从mysql中查找用户信息
		Models.GetUserByAccountNum(db, user, accountNum)
		if user.ID == 0 {
			return errors.New("user not found")
		}
		user.Repair()
		args := redis.Args{accountNum, 3600}.AddFlat(user)

		c := Redis.RedisPoolGet()
		defer Redis.RedisPoolPut(c)
		_, err = redis2.HmsetWithExpireScript.Do(c, args...)
		if err != nil {
			return err
		}

		return nil
	}
}
