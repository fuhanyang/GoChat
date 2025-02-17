package message

import (
	"User/StatusCode"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"rpc/Client"
	"rpc/Message"
	"rpc/Service"
	"rpc/handler"
	"sync"
)

const ServiceName = "Message"

var HandlerNames = handler.ServiceMethods[ServiceName]
var HandlersRequestMap = map[string]*sync.Pool{
	"SendText":    {New: func() interface{} { return &handler.HandlerWithTime{Type: "SendText"} }},
	"RefreshText": {New: func() interface{} { return &Message.RefreshRequest{} }},
}
var HandlersResponseMap = map[string]*sync.Pool{
	"SendText":    {New: func() interface{} { return &Message.SendTextResponse{} }},
	"RefreshText": {New: func() interface{} { return &Message.RefreshTextResponse{} }},
}

type MessageServiceHandler struct {
	Client Message.MessageServiceClient
}

// Init 初始化，注入客户端工厂
func Init() {
	fmt.Println("message service Init")
	var c = MessageServiceHandler{}
	c.InjectClientFactory()

}
func (c *MessageServiceHandler) InjectClientFactory() {
	Service.Inject(ServiceName)
	for _, HandlerName := range HandlerNames {
		handler.InjectHandlers(ServiceName, HandlerName, HandlersRequestMap[HandlerName], HandlersResponseMap[HandlerName])
	}
	//复制一份
	_c := *c
	var factory = func(conn *grpc.ClientConn) Client.Client {
		_c.Client = Message.NewMessageServiceClient(conn)
		return &_c
	}
	Client.Inject(ServiceName, factory)
}

func (c MessageServiceHandler) Handle(ctx context.Context, Handler handler.HandlerRequest) ([]byte, error) {
	var response handler.HandlerResponse
	var err error
	var responseJson []byte
	// 解析包装的handler
	// 如果是带有时间戳的handler，则取出handler
	if _, ok := Handler.(*handler.HandlerWithTime); ok {
		Handler = Handler.(*handler.HandlerWithTime).Handler
	}
	fmt.Println("parsed handler :", Handler)

	switch Handler.GetHandlerName() {
	case "SendText":
		sendTextRequest := Handler.(*Message.SendTextRequest)
		response, _ = c.Client.SendText(ctx, sendTextRequest)
		responseJson, _ = json.Marshal(response.(*Message.SendTextResponse))
	case "RefreshText":
		refreshTextRequest := Handler.(*Message.RefreshRequest)
		response, _ = c.Client.RefreshText(ctx, refreshTextRequest)
		responseJson, _ = json.Marshal(response.(*Message.RefreshTextResponse))
	}
	if response != nil && response.GetCode() != StatusCode.StatusOK {
		fmt.Println("response code :", response.GetCode())
		err = errors.New(Handler.GetHandlerName() + " get error")
	}
	return responseJson, err
}
