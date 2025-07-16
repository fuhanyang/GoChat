package message

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

const ServiceName = "message"

var once sync.Once
var HandlerNames = method.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"SendText":    {New: func() interface{} { return &SendTextRequest{} }},
	"RefreshText": {New: func() interface{} { return &RefreshRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"SendText":    {New: func() interface{} { return &SendTextResponse{} }},
	"RefreshText": {New: func() interface{} { return &RefreshTextResponse{} }},
}
var reqTypeMap = map[string]reflect.Type{
	"SendText":    reflect.TypeOf((*SendTextRequest)(nil)),
	"RefreshText": reflect.TypeOf((*RefreshRequest)(nil)),
}
var respTypeMap = map[string]reflect.Type{
	"SendText":    reflect.TypeOf((*SendTextResponse)(nil)),
	"RefreshText": reflect.TypeOf((*RefreshTextResponse)(nil)),
}

type MessageServiceHandler struct {
	Client MessageServiceClient
}

// Init 初始化，注入客户端工厂

func Init() {
	inject.Inject(ServiceName, func() {
		fmt.Println("message service Init")
		var c = MessageServiceHandler{}
		c.InjectClientFactory()
	})
}
func (c *MessageServiceHandler) InjectClientFactory() {
	once.Do(func() {
		for _, HandlerName := range HandlerNames {
			method.InjectMethod(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
		}
		//复制一份
		_c := *c
		var factory = func(conn *grpc.ClientConn) Client.Client {
			_c.Client = NewMessageServiceClient(conn)
			return &_c
		}
		Client.Inject(ServiceName, factory)
	})
}

func (c *MessageServiceHandler) Handle(ctx context.Context, Handler method.MethodRequest) ([]byte, error) {
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

	//var response method.MethodResponse
	//var err error
	//var responseJson []byte
	//// 解析包装的handler
	//// 如果是带有时间戳的handler，则取出handler
	//if _, ok := Handler.(*method.MessageWithTime); ok {
	//	Handler = Handler.(*method.MessageWithTime).Handler
	//}
	//fmt.Println("parsed method :", Handler)
	//
	//switch Handler.GetHandlerName() {
	//case "SendText":
	//	sendTextRequest := Handler.(*message.SendTextRequest)
	//	response, _ = c.Client.SendText(ctx, sendTextRequest)
	//	responseJson, _ = json.Marshal(response.(*message.SendTextResponse))
	//case "RefreshText":
	//	refreshTextRequest := Handler.(*message.RefreshRequest)
	//	response, _ = c.Client.RefreshText(ctx, refreshTextRequest)
	//	responseJson, _ = json.Marshal(response.(*message.RefreshTextResponse))
	//}
	//if response != nil && response.GetCode() != StatusCode.StatusOK {
	//	fmt.Println("response code :", response.GetCode())
	//	err = errors.New(Handler.GetHandlerName() + " get error")
	//}
	//return responseJson, err
}
