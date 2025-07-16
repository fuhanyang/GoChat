package authentication

import (
	"common/chain"
	redis2 "common/redis"
	"errors"
	"github.com/jinzhu/gorm"
	"user/DAO/Redis"
)

func UserLogOff(accountNum string, password string, ip string, db *gorm.DB) error {
	if accountNum == "" {
		return errors.New("accountNum is empty")
	}

	ctx := chain.LoadHandlers(chain.ZapLogger, chain.DefaultTimer(), distributeLockHandler(), mysqlHandler(), logOffRedisHandler())
	ctx.Set("ip", ip)
	ctx.Set("accountNum", accountNum)
	ctx.Set("password", password)
	ctx.Set("db", db)

	return ctx.Apply()

}

func logOffRedisHandler() chain.MyHandlerFunc {
	return func(ctx *chain.MyContext) error {
		accountNum, err := chain.GetToType[string](ctx, "accountNum")
		if err != nil {
			return err
		}

		// 更新redis
		_, err = Redis.RedisDo(redis2.SET, "is_online"+accountNum, false)
		if err != nil {
			return err
		}
		//向全局的set中添加或删除用户表示是否在线
		if _, err = Redis.RedisDo(redis2.SREM, redis2.ONLINE_USER_SET, accountNum); err != nil {
			return err
		}
		return nil
	}
}
