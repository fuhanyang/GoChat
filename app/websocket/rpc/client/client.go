package client

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	resolver2 "google.golang.org/grpc/resolver"
	"rpc/Client"
	"rpc/Service"
	"rpc/friend"
	"rpc/message"
	"rpc/user"
	"sync"
	"time"
)

var (
	once                 sync.Once
	err                  error
	UserServiceClient    *user.UserServiceHandler
	FriendServiceClient  *friend.FriendServiceHandler
	MessageServiceClient *message.MessageServiceHandler
)

var etcdResolverBuilder resolver2.Builder

func Init(_host string, _port string) {
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
		go initFriendServiceClient()
		go initMessageServiceClient()
	})
}
func initUserServiceClient() {
	var (
		client      Client.Client
		ServiceName = "user"
	)
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
func initFriendServiceClient() {
	var (
		client      Client.Client
		ServiceName = "friend"
	)
	for {
		client, err = Service.ConnService(ServiceName, etcdResolverBuilder)
		if err == nil {
			break
		}
		fmt.Println(err)
		// 每秒重试一次
		time.Sleep(time.Second)
	}
	FriendServiceClient = client.(*friend.FriendServiceHandler)
}
func initMessageServiceClient() {
	var (
		client      Client.Client
		ServiceName = "message"
	)
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
