package message

import (
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"rpc/Service/method"
)

type MessageWithTime struct {
	Handler   method.MethodRequest `json:"method"`
	CreatedAt string               `json:"created_at"`
	Type      string               `json:"type"`
}

func (h *MessageWithTime) ProtoMessage() {
	h.Handler.ProtoMessage()
}
func (h *MessageWithTime) Reset() {
	h.Handler.Reset()
}
func (h *MessageWithTime) String() string {
	return h.Handler.String()
}
func (h *MessageWithTime) ProtoReflect() protoreflect.Message {
	return h.Handler.ProtoReflect()
}
func (h *MessageWithTime) Descriptor() ([]byte, []int) {
	return h.Handler.Descriptor()
}
func (h *MessageWithTime) GetHandlerName() string {
	return h.Handler.GetHandlerName()
}

// UnmarshalJSON  实现自定义的反序列化方法，能自动识别Handler接口的类型，并反序列化
func (h *MessageWithTime) UnmarshalJSON(data []byte) error {
	var err error
	//这里不定义新类型在下面调用json.Unmarshal(data, &aux)时会递归执行，导致崩溃
	type Temp MessageWithTime
	aux := &struct {
		Raw json.RawMessage `json:"method"`
		*Temp
	}{
		Temp: (*Temp)(h),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch aux.Type {
	case "SendText":
		var handler SendTextRequest
		if err = json.Unmarshal(aux.Raw, &handler); err != nil {
			return err
		}
		handler.Time = []byte(aux.CreatedAt)
		h.Handler = &handler
	default:
		return errors.New("unknown type")
	}

	return nil
}
