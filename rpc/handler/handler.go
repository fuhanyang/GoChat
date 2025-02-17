package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"sync"
)

var ServiceMethodRequestMap = make(map[string]map[string]*sync.Pool)
var ServiceMethodResponseMap = make(map[string]map[string]*sync.Pool)
var ServiceMethods = map[string][]string{
	"User":    {"Register", "Login", "Logoff", "FindUser"},
	"Message": {"SendText", "RefreshText"},
	"Friend":  {"AddFriend", "GetFriends"},
}

func InjectHandlers(serviceName string, handlerName string, reqPool *sync.Pool, respPool *sync.Pool) {
	if ServiceMethodRequestMap[serviceName] == nil {
		ServiceMethodRequestMap[serviceName] = make(map[string]*sync.Pool)
	}
	ServiceMethodRequestMap[serviceName][handlerName] = reqPool
	if ServiceMethodResponseMap[serviceName] == nil {
		ServiceMethodResponseMap[serviceName] = make(map[string]*sync.Pool)
	}
	ServiceMethodResponseMap[serviceName][handlerName] = reqPool
}

type HandlerRequest interface {
	ProtoMessage()
	Reset()
	String() string
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	GetHandlerName() string
}
type HandlerResponse interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	GetHandlerName() string
	GetCode() int32
}

func GetHandlersType(ServiceName string, HandlerName string) HandlerRequest {
	return ServiceMethodRequestMap[ServiceName][HandlerName].Get().(HandlerRequest)
	//if ServiceName == "Message" && HandlerName != "RefreshText" {
	//	return &HandlerWithTime{Type: HandlerName}
	//}
	//switch HandlerName {
	//case "Register":
	//	return &User.RegisterRequest{}
	//case "Login":
	//	return &User.LoginRequest{}
	//case "Logoff":
	//	return &User.LogOffRequest{}
	//case "FindUser":
	//	return &User.FindUserRequest{}
	//case "RefreshText":
	//	return &Message.RefreshRequest{}
	//case "AddFriend":
	//	return &Friend.AddFriendRequest{}
	//case "GetFriends":
	//	return &Friend.GetFriendsRequest{}
	//default:
	//	return nil
	//}
}
