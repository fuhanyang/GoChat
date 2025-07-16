package Logic

import (
	"friend/Models"
	"github.com/jinzhu/gorm"
)

func AddFriend(db *gorm.DB, userAccountNum string, friendAccountNum string, userName string, friendName string) error {
	return Models.CreateFriend(db, userAccountNum, friendAccountNum, userName, friendName)
}
