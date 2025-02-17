package handler

import (
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"rpc/Message"
)

type HandlerWithTime struct {
	Handler   HandlerRequest `json:"handler"`
	CreatedAt string         `json:"created_at"`
	Type      string         `json:"type"`
}

func (h *HandlerWithTime) ProtoMessage() {
	h.Handler.ProtoMessage()
}
func (h *HandlerWithTime) Reset() {
	h.Handler.Reset()
}
func (h *HandlerWithTime) String() string {
	return h.Handler.String()
}
func (h *HandlerWithTime) ProtoReflect() protoreflect.Message {
	return h.Handler.ProtoReflect()
}
func (h *HandlerWithTime) Descriptor() ([]byte, []int) {
	return h.Handler.Descriptor()
}
func (h *HandlerWithTime) GetHandlerName() string {
	return h.Handler.GetHandlerName()
}

// UnmarshalJSON  实现自定义的反序列化方法，能自动识别Handler接口的类型，并反序列化
func (h *HandlerWithTime) UnmarshalJSON(data []byte) error {
	var err error
	//这里不定义新类型在下面调用json.Unmarshal(data, &aux)时会递归执行，导致崩溃
	type Temp HandlerWithTime
	aux := &struct {
		Raw json.RawMessage `json:"handler"`
		*Temp
	}{
		Temp: (*Temp)(h),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch aux.Type {
	case "SendText":
		var handler Message.SendTextRequest
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
