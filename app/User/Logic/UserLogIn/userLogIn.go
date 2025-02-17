package UserLogIn

import (
	"User/Logic/UserStatusChange"
)

func UserLogIn(accountNum string, password string, ip string) error {
	err := UserStatusChange.ChangeUserOnlineStatus(ip, password, accountNum, true)
	return err
}
