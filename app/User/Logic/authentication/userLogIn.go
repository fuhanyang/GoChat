package authentication

import (
	"common/chain"
	redis2 "common/redis"
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"log"
	Websocket "rpc/websocket"
	"user/DAO/Redis"
	"user/Logic/Search/SearchUser"
	"user/Logic/UserCreate"
	"user/Models"
	"user/rpc/client"
)

func UserLogIn(accountNum string, password string, ip string, db *gorm.DB) error {
	if accountNum == "" {
		return errors.New("accountNum is empty")
	}

	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), distributeLockHandler(), mysqlHandler(), redisHandler())
	ctx.Set("ip", ip)
	ctx.Set("accountNum", accountNum)
	ctx.Set("password", password)
	ctx.Set("db", db)

	return ctx.Apply()
}

func distributeLockHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		ip, err := chain.GetToType[string](ctx, "ip")
		if err != nil {
			return err
		}
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}
		password, err := chain.GetToType[string](ctx, "password")
		if err != nil {
			return err
		}

		user, err := UserCreate.CreateUser("", ip, password, accountNum)
		if err != nil {
			return err
		}
		//对此账号的用户加上分布式锁
		// 采用阻塞锁
		err = user.RedisBlockLock()
		if err != nil {
			fmt.Println("服务器繁忙，请稍后再试")
			return err
		}
		defer func() {
			// 释放资源
			err = user.RedisUnlock()
			if err != nil {
				panic(err)
			}
			Models.ReleaseUser(user)
		}()
		ctx.Set("user", user)
		// 执行后续中间件
		ctx.Next()
		return nil
	}
}

func mysqlHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}
		db, err := chain.GetToType[*gorm.DB](ctx, "db")
		if err != nil {
			return err
		}
		user, err := chain.GetToType[*Models.User](ctx, "user")
		if err != nil {
			return err
		}

		// 从数据库中查询用户信息
		_user, err := SearchUser.SearchUserByAccountNum(accountNum, db)

		if err != nil {
			return err
		}
		//校验密码是否正确
		if user.GetPassword() != _user.GetPassword() {
			return errors.New("Incorrect password")
		}
		return nil
	}
}

func redisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}

		res, err := Redis.RedisDo(redis2.GET, "is_online"+accountNum)
		status, err := redis.Bool(res, err)
		if errors.Is(err, redis.ErrNil) {
			status = false
		} else if err != nil {
			return err
		}
		if !status {
			// 更新redis
			_, err = Redis.RedisDo(redis2.SET, "is_online"+accountNum, true)
			if err != nil {
				return err
			}
			//向全局的set中添加或删除用户表示是否在线
			if _, err = Redis.RedisDo(redis2.SADD, redis2.ONLINE_USER_SET, accountNum); err != nil {
				return err
			}
		} else {
			// 已经在线，踢掉之前的连接
			_, err = client.WebsocketServiceClient.Client.KickUser(context.Background(), &Websocket.KickUserRequest{
				AccountNum:  accountNum,
				Reason:      "已在其他地方登录，您已被迫下线",
				HandlerName: "KickUser",
			})
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	}
}
