package client

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	resolver2 "google.golang.org/grpc/resolver"
	"rpc/Client"
	"rpc/Service"
	"rpc/websocket"
	"sync"
	"time"
)

var (
	once                   sync.Once
	err                    error
	WebsocketServiceClient *websocket.WebsocketServiceHandler
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
		go initWebsocketServiceClient()

	})
}
func initWebsocketServiceClient() {
	var (
		client      Client.Client
		ServiceName = "websocket"
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
	WebsocketServiceClient = client.(*websocket.WebsocketServiceHandler)
}
