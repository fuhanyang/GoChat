package SearchUser

import (
	Redis2 "User/DAO/Redis"
	"User/Models"
	"errors"
	"fmt"
	redis2 "github.com/gomodule/redigo/redis"
)

// SearchUserByAccountNum 在数据库中通过账号查找用户信息
func SearchUserByAccountNum(accountNum string) (*Models.User, error) {
	user := Models.NewUser()
	// 判断在redis中是否存在该用户信息
	//先从redis中查找用户信息，如果没有再从mysql中查找
	var values []interface{}
	reply, err := Redis2.RedisDo(Redis2.HGETALL, accountNum)
	if err != nil && !errors.Is(err, redis2.ErrNil) {
		//产生了其他错误则要返回
		return user, err
	}
	values, err = redis2.Values(reply, err)
	//redis中不存在该用户信息则从mysql中查找
	if errors.Is(err, redis2.ErrNil) || len(values) == 0 {
		goto mysql
	}
	err = redis2.ScanStruct(values, user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("redis中获取用户信息成功")
	return user, err
mysql:
	//从mysql中查找用户信息
	Models.GetUserByAccountNum(user, accountNum)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	user.Repair()
	args := redis2.Args{accountNum}.AddFlat(user)
	_, err = Redis2.RedisDo(Redis2.HMSET, args...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("mysql中获取用户信息成功")
	return user, nil
}
