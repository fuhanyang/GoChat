package Logic

import (
	"Friend/Models"
	"rpc/Friend"
)

var limit = 10

func GetFriends(accountNum string) []*Friend.Friend {
	friends := Models.GetFriend(accountNum, limit)
	_friends := make([]*Friend.Friend, 0, len(friends))
	var friendAccountNum string
	var friendName string
	for _, friend := range friends {
		if friend.User1 != accountNum {
			friendAccountNum = friend.User1
			friendName = friend.User1Name
		} else {
			friendAccountNum = friend.User2
			friendName = friend.User2Name
		}
		_friends = append(_friends, &Friend.Friend{
			AccountNum: friendAccountNum,
			Name:       friendName,
		})
	}
	return _friends
}
