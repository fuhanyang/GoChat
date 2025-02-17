package Service

import (
	"Message/Logic/GetMessage"
	"Message/Logic/SendMessage"
	"User/StatusCode"
	"rpc/Message"

	"context"
	"fmt"
	"time"
)

func (s *server) RefreshText(ctx context.Context, req *Message.RefreshRequest) (*Message.RefreshTextResponse, error) {
	var (
		res Message.RefreshTextResponse
	)
	res.SenderAccountNum = req.SenderAccountNum
	res.ReceiverAccountNum = req.ReceiverAccountNum

	contents := GetMessage.RefreshText(req.SenderAccountNum, req.ReceiverAccountNum)
	for _, content := range contents {
		res.Content = append(res.Content, content)
	}
	res.Msg = "Refresh Text Success"
	res.Code = 200
	res.HandlerName = "RefreshText"
	return &res, nil
}
func (s *server) SendText(ctx context.Context, req *Message.SendTextRequest) (*Message.SendTextResponse, error) {
	var (
		res Message.SendTextResponse
		err error
	)
	res.HandlerName = "SendText"
	content := []byte(req.Content)
	timeStr := string(req.Time)
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	fmt.Println(parsedTime)
	if err != nil {
		goto ERR
	}
	err = SendMessage.SendText(req.SenderAccountNum, req.ReceiverAccountNum, int64(len(content)), content, parsedTime)
	if err != nil {
		goto ERR
	}
	res.Msg = "SendText Success"
	res.Code = StatusCode.StatusOK
	return &res, err
ERR:
	fmt.Println(err)
	res.Msg = "Error :" + err.Error()
	res.Code = 500
	return &res, nil

}
