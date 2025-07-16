package websocket

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

const ServiceName = "websocket"

var once sync.Once
var HandlerNames = method.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"TryConnect": {New: func() interface{} { return &TryConnectRequest{} }},
	"KickUser":   {New: func() interface{} { return &KickUserRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"TryConnect": {New: func() interface{} { return &TryConnectResponse{} }},
	"KickUser":   {New: func() interface{} { return &KickUserResponse{} }},
}
var reqTypeMap = map[string]reflect.Type{
	"TryConnect": reflect.TypeOf((*TryConnectRequest)(nil)),
	"KickUser":   reflect.TypeOf((*KickUserRequest)(nil)),
}
var respTypeMap = map[string]reflect.Type{
	"TryConnect": reflect.TypeOf((*TryConnectResponse)(nil)),
	"KickUser":   reflect.TypeOf((*KickUserResponse)(nil)),
}

type WebsocketServiceHandler struct {
	Client WebsocketServiceClient
}

func Init() {
	inject.Inject(ServiceName, func() {
		fmt.Println("websocket service Init")
		var c = WebsocketServiceHandler{}
		c.InjectClientFactory()
	})
}
func (c *WebsocketServiceHandler) Handle(ctx context.Context, Handler method.MethodRequest) ([]byte, error) {
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

}

func (c *WebsocketServiceHandler) InjectClientFactory() {
	once.Do(func() {
		for _, HandlerName := range HandlerNames {
			method.InjectMethod(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
		}
		//复制一份
		_c := *c
		var factory = func(conn *grpc.ClientConn) Client.Client {
			_c.Client = NewWebsocketServiceClient(conn)
			return &_c
		}
		Client.Inject(ServiceName, factory)
	})

}
