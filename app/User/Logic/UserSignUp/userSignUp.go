package UserSignUp

import (
	Redis2 "User/DAO/Redis"
	"User/Logic/UserCreate"
	"User/Models"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// UserSignUp 注册用户
func UserSignUp(username string, password string, ip string) (string, error) {
	// 创建用户
	user, err := UserCreate.CreateUser(username, ip, password, UserCreate.CreateAccountNum())
	defer user.Release()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// TODO: 写入数据库,把账号通过回调队列给server，server再给网关层然后写入bitmap
	Models.WriteUser(user)
	args := redis.Args{user.GetAccountNum()}.AddFlat(user)
	_, err = Redis2.RedisDo(Redis2.HMSET, args...)
	if err != nil {
		return "write redis error", err
	}
	fmt.Println("User Sign Up Success")
	return user.GetAccountNum(), nil
}
