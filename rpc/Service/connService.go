package Service

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"rpc/Client"
)

var ServicesName []string

// Inject 注入rpc服务以及其对应的handler的request和response实例
func Inject(serviceName string) {
	ServicesName = append(ServicesName, serviceName)
}
func ConnService(ServiceName string) (Client.Client, error) {

	// 连接rpc服务
	addr := GetServerAddr(ServiceName)
	if addr == "" {
		return nil, fmt.Errorf("service not found")
	}
	fmt.Println(ServiceName, " service found at ", addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c := Client.NewClient(conn, ServiceName)
	return c, nil
}
