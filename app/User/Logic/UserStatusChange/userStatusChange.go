package UserStatusChange

import (
	"User/DAO/Mysql"
	Redis2 "User/DAO/Redis"
	"User/Logic/Search/SearchUser"
	"User/Logic/UserCreate"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//用worker复用来实现对mq的消费，这里暂时先不实现，不使用http请求，不过为了实现中间件的效果，要进行责任链的封装

func ChangeName(name string, password string, accountNum string) error {
	user, err := UserCreate.CreateUser("", "", password, accountNum)
	if err != nil {
		return err
	}
	if user == nil {
		panic("User is nil")
	}
	// 采用阻塞锁，如果已经被占用则一直等待
	err = user.RedisBlockLock()
	if err != nil {
		return err
	}
	defer func() {
		err = user.RedisUnlock()
		if err != nil {
			panic(err)
		}
		// 释放资源
		user.Release()
	}()
	// 从数据库中查询用户信息
	_user, err := SearchUser.SearchUserByAccountNum(accountNum)
	if err != nil {
		return err
	}
	//校验密码是否正确
	if user.GetPassword() != _user.GetPassword() {
		return errors.New("Incorrect password")
	}
	user.SetName(name)
	//TODO: 更新redis
	args := redis.Args{user.GetAccountNum()}.Add("name", user.GetName())
	_, err = Redis2.RedisDo(Redis2.HSET, args...)
	return err
}

func ChangePassword(password string, accountNum string, newPassword string) error {
	user, err := UserCreate.CreateUser("", "", password, accountNum)
	if err != nil {
		return err
	}
	if user == nil {
		panic("User is nil")
	}
	// 采用阻塞锁，如果已经被占用则一直等待
	err = user.RedisBlockLock()
	if err != nil {
		return err
	}
	defer func() {
		err = user.RedisUnlock()
		if err != nil {
			panic(err)
		}
		// 释放资源
		user.Release()
	}()
	// 从数据库中查询用户信息
	_user, err := SearchUser.SearchUserByAccountNum(accountNum)

	if err != nil {
		return err
	}
	//校验密码是否正确
	if user.GetPassword() != _user.GetPassword() {
		return errors.New("Incorrect password")
	}
	user.SetPassword(newPassword)
	//TODO: 更新redis
	args := redis.Args{user.GetAccountNum()}.Add("password", user.GetPassword())
	_, err = Redis2.RedisDo(Redis2.HSET, args...)
	return err
}

// ChangeUserOnlineStatus  更改用户在线状态
func ChangeUserOnlineStatus(ip string, password string, accountNum string, status bool) error {
	//对此账号的用户加上分布式锁
	user, err := UserCreate.CreateUser("", ip, password, accountNum)
	if err != nil {
		return err
	}
	// 采用非阻塞锁，如果已经被占用则直接退出
	err = user.RedisBlockLock()
	if err != nil {
		fmt.Println(" redis lock error: ", err)
		return err
	}
	defer func() {
		err = user.RedisUnlock()
		if err != nil {
			panic(err)
		}
		// 释放资源
		user.Release()
	}()
	// 从数据库中查询用户信息
	_user, err := SearchUser.SearchUserByAccountNum(accountNum)

	if err != nil {
		return err
	}
	//校验密码是否正确
	if user.GetPassword() != _user.GetPassword() {
		return errors.New("Incorrect password")
	}
	if _user.GetIsOnline() == status {
		var ss string
		if status {
			ss = "online"
		} else {
			ss = "offline"
		}
		return errors.New("Online status is already " + ss)
	}
	user.SetIsOnline(status)
	// 更新数据库
	args := redis.Args{accountNum}.AddFlat(user)
	_, err = Redis2.RedisDo(Redis2.HMSET, args...)
	if err != nil {
		fmt.Println(err)
	}
	Mysql.MysqlDb.Model(user).Where("account_num = ?", accountNum).Update("is_online", status)
	return nil
}

// CheckOnlineStatus  检查用户在线状态
func CheckOnlineStatus(password string, accountNum string, expect bool) error {
	user, err := UserCreate.CreateUser("", "", password, accountNum)
	if err != nil {
		return err
	}
	if user == nil {
		panic("User is nil")
	}
	// 采用阻塞锁，如果已经被占用则一直等待
	err = user.RedisBlockLock()
	if err != nil {
		return err
	}
	defer func() {
		err = user.RedisUnlock()
		if err != nil {
			panic(err)
		}
		// 释放资源
		user.Release()
	}()
	// 从数据库中查询用户信息
	_user, err := SearchUser.SearchUserByAccountNum(accountNum)

	if err != nil {
		return err
	}
	//校验密码是否正确
	if user.GetPassword() != _user.GetPassword() {
		return errors.New("Incorrect password")
	}
	var ToString func(bool) string
	ToString = func(expect bool) string {
		if expect {
			return "online"
		}
		return "offline"
	}
	if user.GetIsOnline() != expect {
		return errors.New("Online status is wrong, expect " + ToString(expect))
	}
	return nil
}
