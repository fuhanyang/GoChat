package client

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	resolver2 "google.golang.org/grpc/resolver"
	"rpc/Client"
	"rpc/Service"
	"rpc/message"
	"rpc/user"
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
var etcdResolverBuilder resolver2.Builder

func Init(_host string, _port string) {
	host = _host
	port = _port
	addr := fmt.Sprintf("http://%s:%s", _host, _port)
	etcdClient, err := clientv3.NewFromURL(addr)
	fmt.Println("connect etcd success at ", addr)
	if err != nil {
		panic(err)
	}
	etcdResolverBuilder, err = resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	once.Do(func() {
		go initUserServiceClient()
		go initMessageServiceClient()
	})
}
func initUserServiceClient() {
	var (
		client      Client.Client
		ServiceName = "user"
	)
	//go func() {
	//	err := service.WatchServiceName(ServiceName, host, port)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//time.Sleep(1 * time.Second)
	for {
		client, err = Service.ConnService(ServiceName, etcdResolverBuilder)
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
		ServiceName = "message"
	)
	//go func() {
	//	err := service.WatchServiceName(ServiceName, host, port)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//time.Sleep(1 * time.Second)
	for {
		client, err = Service.ConnService(ServiceName, etcdResolverBuilder)
		if err == nil {
			break
		}
		fmt.Println(err)
		// 每秒重试一次
		time.Sleep(time.Second)
	}
	MessageServiceClient = client.(*message.MessageServiceHandler)
}
