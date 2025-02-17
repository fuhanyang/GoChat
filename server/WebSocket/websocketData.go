package WebSocket

import (
	"time"
)

type WebSocketData struct {
	Data       string `json:"data"`
	AccountNum string `json:"accountNum"`
	Receiver   string `json:"receiver"`
	Type       string `json:"type"`
	CreatedAt  string `json:"createdAt"`
}

func FormErrWebSocketData(err error, accountNum string, receiver string) WebSocketData {
	return WebSocketData{
		Data:       err.Error(),
		AccountNum: accountNum,
		Receiver:   receiver,
		Type:       "error",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
}
