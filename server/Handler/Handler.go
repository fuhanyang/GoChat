package Handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"rpc/Message"
)

type Handler interface {
	ProtoMessage()
	Reset()
	String() string
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
	GetHandlerName() string
}
type MessageHandler interface {
	Handler
	GetContent() string
}

// GetHandlersType 返回指定Handler的实例，包含发送消息的一系列操作
func GetHandlersType(HandlerName string) MessageHandler {
	switch HandlerName {
	case "SendText":
		return &Message.SendTextRequest{}
	}
	return nil
}
