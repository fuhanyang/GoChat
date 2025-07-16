package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"rpc/Service"
	Websocket "rpc/websocket"
	settings "websocket/settings"
)

type server struct {
	Websocket.UnimplementedWebsocketServiceServer
}

var service = &Service.Service{}

func InitServer(config *settings.ServiceConfig) {
	service.Host = config.Host
	service.Port = config.Port
	service.Name = config.Name
	service.Protocol = config.Protocol
	fmt.Println("service Init Success, Info: ", service)
}
func NewServer(ctx context.Context, EtcdHost string, EtcdPort string) *grpc.Server {
	s := grpc.NewServer()
	Websocket.RegisterWebsocketServiceServer(s, &server{})
	// 服务注册
	err := Service.ServiceRegister(service, ctx, EtcdHost, EtcdPort)
	if err != nil {
		panic(err)
	}
	fmt.Println("service Register Success, Info: ", service)
	return s
}
