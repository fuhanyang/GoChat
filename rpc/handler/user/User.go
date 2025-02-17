package user

import (
	"User/StatusCode"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"rpc/Client"
	"rpc/Service"
	"rpc/User"
	"rpc/handler"
	"sync"
)

const ServiceName = "User"

var HandlerNames = handler.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"Register": {New: func() interface{} { return &User.RegisterRequest{} }},
	"Login":    {New: func() interface{} { return &User.LoginRequest{} }},
	"Logoff":   {New: func() interface{} { return &User.LogOffRequest{} }},
	"FindUser": {New: func() interface{} { return &User.FindUserRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"Register": {New: func() interface{} { return &User.RegisterResponse{} }},
	"Login":    {New: func() interface{} { return &User.LoginResponse{} }},
	"Logoff":   {New: func() interface{} { return &User.LogOffRequest{} }},
	"FindUser": {New: func() interface{} { return &User.FindUserRequest{} }},
}

type UserServiceHandler struct {
	Client User.UserServiceClient
}

// Init 初始化，注入客户端工厂和服务端请求与响应
func Init() {
	fmt.Println("user service Init")
	var c = UserServiceHandler{}
	c.InjectClientFactory()

}
func (c *UserServiceHandler) InjectClientFactory() {
	Service.Inject(ServiceName)
	for _, HandlerName := range HandlerNames {
		handler.InjectHandlers(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
	}
	//复制一份
	_c := *c
	var factory = func(conn *grpc.ClientConn) Client.Client {
		_c.Client = User.NewUserServiceClient(conn)
		return &_c
	}
	Client.Inject(ServiceName, factory)
}
func (c UserServiceHandler) Handle(ctx context.Context, Handler handler.HandlerRequest) ([]byte, error) {
	var response handler.HandlerResponse
	var err error
	var responseJson []byte
	switch Handler.GetHandlerName() {
	case "Register":
		registerRequest := Handler.(*User.RegisterRequest)
		response, err = c.Client.Register(ctx, registerRequest)
		if err != nil {
			fmt.Println("user service register error:", err)
		}
		responseJson, _ = json.Marshal(response.(*User.RegisterResponse))
	case "Login":
		loginRequest := Handler.(*User.LoginRequest)
		response, _ = c.Client.Login(ctx, loginRequest)
		responseJson, _ = json.Marshal(response.(*User.LoginResponse))
	case "Logoff":
		logOffRequest := Handler.(*User.LogOffRequest)
		response, _ = c.Client.LogOff(ctx, logOffRequest)
		responseJson, _ = json.Marshal(response.(*User.LogOffResponse))
	case "FindUser":
		findUserRequest := Handler.(*User.FindUserRequest)
		response, _ = c.Client.FindUser(ctx, findUserRequest)
		responseJson, _ = json.Marshal(response.(*User.FindUserResponse))
	}
	if response != nil && response.GetCode() != StatusCode.StatusOK {
		err = errors.New(Handler.GetHandlerName() + " get error")
	}
	return responseJson, err
}
