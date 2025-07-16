package forcedOffline

import (
	"common/redis"
	"fmt"
	Redis2 "user/DAO/Redis"
)

func ForcedOffline(accountNum string) error {
	// 更新数据库
	_, err := Redis2.RedisDo(redis.SET, "is_online"+accountNum, false)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
