package friend

import (
	"User/StatusCode"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"rpc/Client"
	"rpc/Friend"
	"rpc/Service"
	"rpc/handler"
	"sync"
)

const ServiceName = "Friend"

var HandlerNames = handler.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"GetFriends": {New: func() interface{} { return &Friend.GetFriendsRequest{} }},
	"AddFriend":  {New: func() interface{} { return &Friend.AddFriendRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"GetFriends": {New: func() interface{} { return &Friend.GetFriendsResponse{} }},
	"AddFriend":  {New: func() interface{} { return &Friend.AddFriendRequest{} }},
}

type FriendServiceHandler struct {
	Client Friend.FriendServiceClient
}

// Init 初始化，注入客户端工厂
func Init() {
	fmt.Println("friend service Init")
	var c = FriendServiceHandler{}
	c.InjectClientFactory()

}
func (c *FriendServiceHandler) InjectClientFactory() {
	Service.Inject(ServiceName)
	for _, HandlerName := range HandlerNames {
		handler.InjectHandlers(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
	}
	//复制一份
	_c := *c
	var factory = func(conn *grpc.ClientConn) Client.Client {
		_c.Client = Friend.NewFriendServiceClient(conn)
		return &_c
	}
	Client.Inject(ServiceName, factory)
}
func (c *FriendServiceHandler) Handle(ctx context.Context, Handler handler.HandlerRequest) ([]byte, error) {
	var response handler.HandlerResponse
	var err error
	var responseJson []byte
	switch Handler.GetHandlerName() {
	case "GetFriends":
		getFriendRequest := Handler.(*Friend.GetFriendsRequest)
		response, _ = c.Client.GetFriends(ctx, getFriendRequest)
		responseJson, err = json.Marshal(response.(*Friend.GetFriendsResponse))
	case "AddFriend":
		addFriendRequest := Handler.(*Friend.AddFriendRequest)
		response, _ = c.Client.AddFriend(ctx, addFriendRequest)
		responseJson, err = json.Marshal(response.(*Friend.AddFriendResponse))
	}
	if response != nil && response.GetCode() != StatusCode.StatusOK {
		err = errors.New(Handler.GetHandlerName() + " get error")
	}
	return responseJson, err
}
