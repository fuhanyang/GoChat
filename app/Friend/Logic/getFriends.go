package Logic

import (
	"friend/Models"
	"github.com/jinzhu/gorm"
	"rpc/friend"
)

var limit = 10

func GetFriends(db *gorm.DB, accountNum string) []*friend.Friend {
	friends := Models.GetFriend(db, accountNum, limit)
	_friends := make([]*friend.Friend, 0, len(friends))
	var friendAccountNum string
	var friendName string
	for _, _friend := range friends {
		if _friend.User1 != accountNum {
			friendAccountNum = _friend.User1
			friendName = _friend.User1Name
		} else {
			friendAccountNum = _friend.User2
			friendName = _friend.User2Name
		}
		_friends = append(_friends, &friend.Friend{
			AccountNum: friendAccountNum,
			Name:       friendName,
		})
	}
	return _friends
}
