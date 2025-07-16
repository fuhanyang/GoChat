package Init

import (
	"rpc/Service/inject"
	"rpc/friend"
	"rpc/message"
	"rpc/user"
	"rpc/websocket"
)

func InitService() {
	message.Init()
	user.Init()
	friend.Init()
	websocket.Init()
	for _, initFunc := range inject.InitFunctions {
		initFunc()
	}
}
