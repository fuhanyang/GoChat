package Service

import (
	settings "User/Settings"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"rpc/Service"
	"rpc/User"
)

type server struct {
	User.UnimplementedUserServiceServer
}

var service = &Service.Service{}

func InitServer(config *settings.ServiceConfig) {
	service.Host = config.Host
	service.Port = config.Port
	service.Name = config.Name
	service.Protocol = config.Protocol
	fmt.Println("Service Init Success, Info: ", service)
}
func NewServer(ctx context.Context, EtcdHost string, EtcdPort string) *grpc.Server {
	s := grpc.NewServer()
	User.RegisterUserServiceServer(s, &server{})
	// 服务注册
	err := Service.ServiceRegister(service, ctx, EtcdHost, EtcdPort)
	if err != nil {
		panic(err)
	}
	return s
}
