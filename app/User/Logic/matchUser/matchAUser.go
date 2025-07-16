package MatchAUser

import (
	"common/chain"
	redis2 "common/redis"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"log"
	"user/DAO/Redis"
	"user/Models"
)

// MatchAUser 匹配用户 choice为true表示从mysql中匹配，false表示从缓存中匹配
func MatchAUser(accountNum string, choice bool, db *gorm.DB) (*Models.User, error) {
	var (
		user = new(Models.User)
		ctx  *chain.MyContext
	)
	if choice {
		ctx = chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), mysqlHandler())
	} else {
		ctx = chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), redisHandler(), mysqlHandler())
	}
	ctx.Set("choice", choice)
	ctx.Set("accountNum", accountNum)
	ctx.Set("db", db)
	ctx.Set("user", user)
	err := ctx.Apply()

	return user, err
}

func redisHandler() chain.MyHandlerFunc {
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

		//这里可以先从缓存中优先匹配在线用户，如果缓存中没有，则从数据库中获取
		for reTry := 0; reTry < 4; reTry++ {
			reply, err := Redis.RedisDo(redis2.SRANDMEMBER, redis2.ONLINE_USER_SET, 1)
			_accountNum, err := redis.Strings(reply, err)
			if err != nil || len(_accountNum) == 0 {
				log.Println(err)
				return nil
			}
			if _accountNum[0] != accountNum {
				Models.GetUserByAccountNum(db, user, _accountNum[0])
				if user.AccountNum == "" {
					// 缓存中匹配到的用户不在数据库中，重新匹配
					continue
				}
				ctx.Abort()
				return nil
			}
		}
		return nil
	}

}

func mysqlHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		fmt.Println("没有在线用户，在mysql中查找")
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
		for i := 0; i < 10; i++ {
			err = Models.MatchAUser(db, user)
			if err != nil {
				return errors.New("user not found")
			}
			if user.AccountNum != accountNum {
				fmt.Println("Matching AccountNum Number: ", user.AccountNum)
				return nil
			}
		}
		return errors.New("Retried 10 times, still not found user")
	}
}
