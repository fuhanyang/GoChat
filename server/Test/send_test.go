package Test

import (
	"encoding/json"
	"rpc/Message"
	"server/Mq"
	"server/WebSocket"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	err := Mq.NewConnCh()
	if err != nil {
		t.Error(err)
	}
	defer Mq.ConnClose()
	data, err := json.Marshal(Message.SendTextRequest{
		SenderAccountNum:   "123",
		ReceiverAccountNum: "456",
		Msg:                "msg",
		Code:               200,
		HandlerName:        "SendText",
		Content:            "send text test",
		Time:               nil,
	})
	if err != nil {
		t.Error(err)
	}
	var ws = WebSocket.WebSocketData{
		Data:       string(data),
		AccountNum: "123",
		Receiver:   "456",
		Type:       "SendText",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	err = WebSocket.SendMessage(ws)
	if err != nil {
		t.Error(err)
	}
}
