package UserCreate

import (
	"common/snowflake"
	"strconv"
	"user/Models"
)

func CreateAccountNum() string {
	accStr := strconv.FormatInt(snowflake.GetID(), 10) //生成唯一的账号号
	return accStr
}

// CreateUser 创建用户
func CreateUser(username string, IP string, password string, accountNum string) (*Models.User, error) {
	user := Models.NewUser()
	//装配属性
	user.SetAccountNum(accountNum)
	user.SetPassword(password)
	user.SetName(username)
	user.SetIP(IP)
	user.SetType("user")
	//修补数据,生成锁
	user.Repair()
	return user, nil
}
