package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"reflect"
	"rpc/Client"
	"rpc/Service/inject"
	"rpc/Service/method"
	"sync"
)

const ServiceName = "user"

var once sync.Once
var HandlerNames = method.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"Register":      {New: func() interface{} { return &RegisterRequest{} }},
	"Login":         {New: func() interface{} { return &LoginRequest{} }},
	"LogOff":        {New: func() interface{} { return &LogOffRequest{} }},
	"FindUser":      {New: func() interface{} { return &FindUserRequest{} }},
	"ForcedOffline": {New: func() interface{} { return &ForcedOfflineRequest{} }},
	"GetUserInfo":   {New: func() interface{} { return &GetUserInfoRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"Register":      {New: func() interface{} { return &RegisterResponse{} }},
	"Login":         {New: func() interface{} { return &LoginResponse{} }},
	"LogOff":        {New: func() interface{} { return &LogOffResponse{} }},
	"FindUser":      {New: func() interface{} { return &FindUserResponse{} }},
	"ForcedOffline": {New: func() interface{} { return &ForcedOfflineResponse{} }},
	"GetUserInfo":   {New: func() interface{} { return &GetUserInfoResponse{} }},
}
var reqTypeMap = map[string]reflect.Type{
	"Register":      reflect.TypeOf((*RegisterRequest)(nil)),
	"Login":         reflect.TypeOf((*LoginRequest)(nil)),
	"LogOff":        reflect.TypeOf((*LogOffRequest)(nil)),
	"FindUser":      reflect.TypeOf((*FindUserRequest)(nil)),
	"ForcedOffline": reflect.TypeOf((*ForcedOfflineRequest)(nil)),
	"GetUserInfo":   reflect.TypeOf((*GetUserInfoRequest)(nil)),
}
var respTypeMap = map[string]reflect.Type{
	"Register":      reflect.TypeOf((*RegisterResponse)(nil)),
	"Login":         reflect.TypeOf((*LoginResponse)(nil)),
	"LogOff":        reflect.TypeOf((*LogOffResponse)(nil)),
	"FindUser":      reflect.TypeOf((*FindUserResponse)(nil)),
	"ForcedOffline": reflect.TypeOf((*ForcedOfflineResponse)(nil)),
	"GetUserInfo":   reflect.TypeOf((*GetUserInfoResponse)(nil)),
}

type UserServiceHandler struct {
	Client UserServiceClient
}

//  Init 初始化，注入客户端工厂和服务端请求与响应

func Init() {
	inject.Inject(ServiceName, func() {
		fmt.Println("user service Init")
		var c = UserServiceHandler{}
		c.InjectClientFactory()
	})
}
func (c *UserServiceHandler) InjectClientFactory() {
	once.Do(func() {
		for _, HandlerName := range HandlerNames {
			method.InjectMethod(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
		}
		//复制一份
		_c := *c
		var factory = func(conn *grpc.ClientConn) Client.Client {
			_c.Client = NewUserServiceClient(conn)
			return &_c
		}
		Client.Inject(ServiceName, factory)
	})

}
func (c *UserServiceHandler) Handle(ctx context.Context, Handler method.MethodRequest) ([]byte, error) {
	var response method.MethodResponse
	var err error
	var responseJson []byte

	reqType, exists := reqTypeMap[Handler.GetHandlerName()]
	if !exists {
		return nil, errors.New("handler not found")
	}
	respType, exists := respTypeMap[Handler.GetHandlerName()]
	if !exists {
		return nil, errors.New("handler not found")
	}

	handlerValue := reflect.ValueOf(Handler)

	// 如果handler是接口，需要解引用获取具体值
	if handlerValue.Kind() == reflect.Interface && !handlerValue.IsNil() {
		handlerValue = handlerValue.Elem()
	}
	if !handlerValue.IsValid() {
		return nil, errors.New("handler is not valid")
	}
	// 判断类型
	if !handlerValue.Type().AssignableTo(reqType) {
		return nil, errors.New("handler is not assignable")
	}
	// 构造请求值
	reqValue := reflect.New(reqType).Elem()
	reqValue.Set(handlerValue)

	//解析方法
	args := []reflect.Value{reflect.ValueOf(ctx), reqValue}
	v := reflect.ValueOf(c.Client)
	//// 打印可用方法列表，辅助调试
	//fmt.Println("可用方法列表:")
	//for i := 0; i < v.Type().NumMethod(); i++ {
	//	m := v.Type().Method(i)
	//	fmt.Printf("- %s (参数: %d, 返回值: %d)\n",
	//		m.Name,
	//		m.Type.NumIn(),
	//		m.Type.NumOut())
	//}
	Method := v.MethodByName(Handler.GetHandlerName())
	//调用方法
	if !Method.IsValid() {
		return nil, errors.New("method is not valid")
	}
	respValue := Method.Call(args)
	if len(respValue) != 2 {
		return nil, errors.New("handler return value error")
	}
	if respValue[0].Type() != respType {
		return nil, errors.New("handler response type error")
	}
	if !respValue[0].IsNil() {
		response = respValue[0].Interface().(method.MethodResponse)
		responseJson, _ = json.Marshal(response)
	}
	if !respValue[1].IsNil() {
		err = respValue[1].Interface().(error)
	}
	return responseJson, err

	//switch Handler.GetHandlerName() {
	//case "Register":
	//	registerRequest := Handler.(*user.RegisterRequest)
	//	response, err = c.Client.Register(ctx, registerRequest)
	//	if err != nil {
	//		fmt.Println("user service register error:", err)
	//	}
	//	responseJson, _ = json.Marshal(response.(*user.RegisterResponse))
	//case "Login":
	//	loginRequest := Handler.(*user.LoginRequest)
	//	response, _ = c.Client.Login(ctx, loginRequest)
	//	responseJson, _ = json.Marshal(response.(*user.LoginResponse))
	//case "LogOff":
	//	logOffRequest := Handler.(*user.LogOffRequest)
	//	response, _ = c.Client.LogOff(ctx, logOffRequest)
	//	responseJson, _ = json.Marshal(response.(*user.LogOffResponse))
	//case "FindUser":
	//	findUserRequest := Handler.(*user.FindUserRequest)
	//	response, _ = c.Client.FindUser(ctx, findUserRequest)
	//	responseJson, _ = json.Marshal(response.(*user.FindUserResponse))
	//case "ForcedOffline":
	//	forcedOfflineRequest := Handler.(*user.ForcedOfflineRequest)
	//	response, _ = c.Client.ForcedOffline(ctx, forcedOfflineRequest)
	//	responseJson, _ = json.Marshal(response.(*user.ForcedOfflineResponse))
	//}
	//if response != nil && response.GetCode() != StatusCode.StatusOK {
	//	err = errors.New(Handler.GetHandlerName() + " get error")
	//}
	//return responseJson, err
}
