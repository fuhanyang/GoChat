package UserCreate

import (
	"User/Logic/Snowflake"
	"User/Models"
	"strconv"
)

func CreateAccountNum() string {
	accStr := strconv.FormatInt(Snowflake.GetID(), 10) //生成唯一的账号号
	return accStr
}

// CreateUser 创建用户
func CreateUser(username string, IP string, password string, accountNum string) (*Models.User, error) {
	user := Models.NewUser()
	//装配属性
	user.SetAccountNum(accountNum)
	user.SetPassword(password)
	user.SetName(username)
	user.SetIsOnline(false)
	user.SetIP(IP)
	user.SetType("User")
	//修补数据,生成锁
	user.Repair()
	return user, nil
}
