package Logic

import "Friend/Models"

func AddFriend(userAccountNum string, friendAccountNum string, userName string, friendName string) error {
	return Models.CreateFriend(userAccountNum, friendAccountNum, userName, friendName)
}
