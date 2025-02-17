package Init

import (
	"rpc/handler/friend"
	"rpc/handler/message"
	"rpc/handler/user"
)

func InitService() {
	message.Init()
	user.Init()
	friend.Init()
}
