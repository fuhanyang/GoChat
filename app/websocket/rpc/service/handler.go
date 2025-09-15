package service

import (
	"common/jwt"
	"context"
	"errors"
	"rpc/user"
	Websocket "rpc/websocket"
	"websocket/Logic"
	"websocket/Logic/websocket"
	"websocket/rpc/client"
)

var (
	URL = Logic.MachineURL
)

func (s *server) TryConnect(ctx context.Context, req *Websocket.TryConnectRequest) (*Websocket.TryConnectResponse, error) {
	var (
		err   error
		res   = &Websocket.TryConnectResponse{}
		token string
	)
	// 检验用户是否存在
	_user, err := client.UserServiceClient.Client.FindUser(ctx, &user.FindUserRequest{AccountNum: req.GetAccountNum()})
	if err != nil {
		return res, err
	}
	if _user == nil {
		err = errors.New("用户不存在")
		return res, err
	}
	//生成新的jwt token
	token, err = jwt.GenToken(req.GetPassword(), req.GetAccountNum())
	if err != nil {
		return nil, err
	}
	res.Token = token
	res.Url = Logic.MachineURL
	return res, nil
}

func (s *server) KickUser(ctx context.Context, req *Websocket.KickUserRequest) (*Websocket.KickUserResponse, error) {
	var (
		err error
		res = &Websocket.KickUserResponse{}
	)
	res.HandlerName = req.GetHandlerName()

	err = websocket.KickUser(req.Reason, req.GetAccountNum())
	if err != nil {
		res.Code = 500
		res.Msg = "踢下线失败"
		return res, nil
	}
	res.Msg = "用户已被踢下线"
	res.Code = 200
	return res, nil
}
