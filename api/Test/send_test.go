package Test

import (
	"api/Mq"
	WebSocket2 "api/middleware/WebSocket"
	"encoding/json"
	"rpc/message"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	err := Mq.NewConnCh()
	if err != nil {
		t.Error(err)
	}
	defer Mq.ConnClose()
	data, err := json.Marshal(message.SendTextRequest{
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
	var ws = WebSocket2.WebSocketData{
		Data:       string(data),
		AccountNum: "123",
		Receiver:   "456",
		Type:       "SendText",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	err = WebSocket2.SendMessage(ws)
	if err != nil {
		t.Error(err)
	}
}
