package method

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"sync"
)

var ServiceMethodRequestMap = make(map[string]map[string]*sync.Pool)
var ServiceMethodResponseMap = make(map[string]map[string]*sync.Pool)
var ServiceMethods = map[string][]string{
	"user":      {"Register", "Login", "LogOff", "FindUser", "ForcedOffline", "GetUserInfo"},
	"message":   {"SendText", "RefreshText"},
	"friend":    {"AddFriend", "GetFriends", "AddFriendWithAccountNum"},
	"websocket": {"TryConnect", "KickUser"},
}

func InjectMethod(serviceName string, methodName string, reqPool *sync.Pool, respPool *sync.Pool) {
	if ServiceMethodRequestMap[serviceName] == nil {
		ServiceMethodRequestMap[serviceName] = make(map[string]*sync.Pool)
	}
	ServiceMethodRequestMap[serviceName][methodName] = reqPool
	if ServiceMethodResponseMap[serviceName] == nil {
		ServiceMethodResponseMap[serviceName] = make(map[string]*sync.Pool)
	}
	ServiceMethodResponseMap[serviceName][methodName] = respPool
}

type MethodRequest interface {
	ProtoMessage()
	Reset()
	String() string
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	GetHandlerName() string
}
type MethodResponse interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	GetHandlerName() string
	GetCode() int32
}

func GetHandlersType(ServiceName string, HandlerName string) MethodRequest {
	if ServiceMethodRequestMap[ServiceName] == nil {
		return nil
	}
	if ServiceMethodRequestMap[ServiceName][HandlerName] == nil {
		return nil
	}
	return ServiceMethodRequestMap[ServiceName][HandlerName].Get().(MethodRequest)
}
