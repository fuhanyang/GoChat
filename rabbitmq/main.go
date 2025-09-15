package main

import (
	"common/chain"
	"common/viper"
	"common/zap"
	"context"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	resolver2 "google.golang.org/grpc/resolver"
	"os"
	"rabbitmq/ParseMessage"
	"rabbitmq/queue"
	"rabbitmq/settings"
	"rpc/Client"
	"rpc/Service"
	"rpc/Service/Init"
	"rpc/Service/inject"
	"rpc/Service/method"
	"time"
)

var ClientMap = make(map[string]Client.Client)
var queueMap = make(map[string]amqp.Queue)
var ServiceHandlerMsgMap = make(map[string]map[string]<-chan amqp.Delivery)
var etcdResolverBuilder resolver2.Builder
var etcdClient *clientv3.Client

func main() {
	var err error
	mode := flag.String("mode", "local", "运行模式，可选值：local(默认)、debug")
	flag.Parse()

	err = viper.Init(settings.Config, *mode)
	if err != nil {
		panic(err)
	}
	fmt.Println("config init success mode:", *mode)

	// 初始化zap日志
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	zap.InitLogger(path)

	Init.InitService()
	// 连接消息队列
	err = queue.NewConnCh(settings.Config.RabbitMQConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect rabbitmq success at", settings.Config.RabbitMQConfig.Host, settings.Config.RabbitMQConfig.Port)
	defer queue.ConnClose()
	// 启动服务发现
	fmt.Println("start service discovery at", settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)
	StartService(settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)
	// 监听消息队列
	ListenQueues()
	select {}
}

// StartService 启动服务发现
func StartService(host string, port string) {
	var c Client.Client
	var err error
	// 服务发现
	for _, ServiceName := range inject.ServicesName {
		for {
			addr := fmt.Sprintf("http://%s:%s", host, port)
			etcdClient, err = clientv3.NewFromURL(addr)
			fmt.Println("connect etcd success at ", addr)
			if err != nil {
				panic(err)
			}
			etcdResolverBuilder, err = resolver.NewBuilder(etcdClient)
			c, err = Service.ConnService(ServiceName, etcdResolverBuilder)
			if err == nil {
				break
			}
			fmt.Println(err.Error())
			time.Sleep(1 * time.Second)
		}

		ClientMap[ServiceName] = c
		// 绑定消息队列
		ServiceHandlerMsgMap[ServiceName] = make(map[string]<-chan amqp.Delivery)
		for _, HandlerName := range method.ServiceMethods[ServiceName] {
			q, err := queue.QueueDeclare(HandlerName, false, false, false, false, nil)
			if err != nil {
				panic(err)
			}
			queueMap[HandlerName] = q
			fmt.Println("bind queue ", q.Name)
			msg, err := queue.GetMessages(
				q.Name, "",
				false, false, false, false,
				nil,
			)
			if err != nil {
				panic(err)
			}
			ServiceHandlerMsgMap[ServiceName][HandlerName] = msg
		}
	}
}

// ListenQueues 监听消息队列
func ListenQueues() {
	// 监听消息队列
	for ServiceName, HandlerMap := range ServiceHandlerMsgMap {
		c := ClientMap[ServiceName]
		for HandlerName, msg := range HandlerMap {
			go func() {
				// 处理消息
				ctx := context.Background()
				for delivery := range msg {
					_handler, err := ParseMessage.ParseDelivery(method.GetHandlersType(ServiceName, HandlerName), delivery)
					if err != nil {
						fmt.Println(err)
						continue
					}
					delivery.Ack(false)
					_handler.GetHandlerName()
					// 调用rpc服务并返回结果
					res, err := c.Handle(ctx, _handler)
					if err != nil {
						chain.ZapLogger("error", "rpc handle error: %s", err.Error())
					}
					if res == nil {
						fmt.Println("res is nil")
						continue
					}
					chain.ZapLogger("info", "rpc handle success: %s", res)
					// 发送响应给server
					err = queue.GivesResponseTo(delivery.ReplyTo, delivery.CorrelationId, res)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}
}
