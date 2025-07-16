package TempUserCreate

import (
	"errors"
	"user/Models"
)

func CreateTempUser(username string, ip string) error {
	user := Models.NewTempUser()
	if user == nil {
		return errors.New("Failed to create new temp user")
	}
	return nil
}
