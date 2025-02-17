package UserLogOff

import (
	"User/Logic/UserStatusChange"
)

func UserLogOff(accountNum string, password string, ip string) error {
	err := UserStatusChange.ChangeUserOnlineStatus(ip, password, accountNum, false)
	return err
}
