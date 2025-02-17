package client

import (
	"fmt"
	"rpc/Client"
	"rpc/Service"
	"rpc/handler/message"
	"rpc/handler/user"
	"sync"
	"time"
)

var (
	once                 sync.Once
	err                  error
	UserServiceClient    *user.UserServiceHandler
	MessageServiceClient *message.MessageServiceHandler
)
var host string
var port string

func Init(_host string, _port string) {
	host = _host
	port = _port
	once.Do(func() {
		go initUserServiceClient()
		go initMessageServiceClient()
	})
}
func initUserServiceClient() {
	var (
		client      Client.Client
		ServiceName = "User"
	)
	go func() {
		err := Service.WatchServiceName(ServiceName, host, port)
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(1 * time.Second)
	for {
		client, err = Service.ConnService(ServiceName)
		if err == nil {
			break
		}
		fmt.Println(err)
		// 每秒重试一次
		time.Sleep(time.Second)
	}
	UserServiceClient = client.(*user.UserServiceHandler)
}
func initMessageServiceClient() {
	var (
		client      Client.Client
		ServiceName = "Message"
	)
	go func() {
		err := Service.WatchServiceName(ServiceName, host, port)
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(1 * time.Second)
	for {
		client, err = Service.ConnService(ServiceName)
		if err == nil {
			break
		}
		fmt.Println(err)
		// 每秒重试一次
		time.Sleep(time.Second)
	}
	MessageServiceClient = client.(*message.MessageServiceHandler)
}
