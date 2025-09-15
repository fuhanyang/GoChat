package service

import (
	"message/DAO/Mysql"
	"message/Logic/GetMessage"
	"message/Logic/SendMessage"
	"rpc/message"
	"user/Const"

	"context"
	"fmt"
	"time"
)

func (s *server) RefreshText(ctx context.Context, req *message.RefreshRequest) (*message.RefreshTextResponse, error) {
	var (
		res message.RefreshTextResponse
	)
	res.SenderAccountNum = req.SenderAccountNum
	res.ReceiverAccountNum = req.ReceiverAccountNum

	contents := GetMessage.RefreshText(Mysql.MysqlDb, req.SenderAccountNum, req.ReceiverAccountNum)
	for _, content := range contents {
		res.Content = append(res.Content, content)
	}
	res.Msg = "Refresh Text Success"
	res.Code = 200
	res.HandlerName = "RefreshText"
	return &res, nil
}
func (s *server) SendText(ctx context.Context, req *message.SendTextRequest) (*message.SendTextResponse, error) {
	var (
		res message.SendTextResponse
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
	err = SendMessage.SendText(Mysql.MysqlDb, req.SenderAccountNum, req.ReceiverAccountNum, int64(len(content)), content, parsedTime)
	if err != nil {
		goto ERR
	}
	res.Msg = "SendText Success"
	res.Code = Const.StatusOK
	return &res, nil
ERR:
	fmt.Println(err)
	res.Msg = "Error :" + err.Error()
	res.Code = 500
	return &res, nil

}
