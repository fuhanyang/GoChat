package MatchAUser

import (
	"User/Models"
	"errors"
	"fmt"
)

func MatchAUser(accountNum string) (*Models.User, error) {
	var user *Models.User
	fmt.Println("User Account Number: ", accountNum)
	for {
		user = Models.MatchAUser()
		if user == nil {
			return nil, errors.New("User not found")
		}
		if user.AccountNum != accountNum {
			fmt.Println("Matching Account Number: ", user.AccountNum)
			return user, nil
		}
	}
}
