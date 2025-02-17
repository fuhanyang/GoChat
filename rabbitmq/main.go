package main

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"rabbitmq/Mq"
	"rabbitmq/ParseMessage"
	"rabbitmq/Settings"
	"rpc/Client"
	"rpc/Service"
	"rpc/handler"
	"rpc/handler/Init"
	"time"
)

var ClientMap = make(map[string]Client.Client)
var queueMap = make(map[string]amqp.Queue)
var ServiceHandlerMsgMap = make(map[string]map[string]<-chan amqp.Delivery)

// StartService 启动服务发现
func StartService(host string, port string) {
	var c Client.Client
	var err error
	// 服务发现
	for _, ServiceName := range Service.ServicesName {
		go func() {
			err := Service.WatchServiceName(ServiceName, host, port)
			if err != nil {
				fmt.Println(err)
			}
		}()
		time.Sleep(1 * time.Second)
		for {
			c, err = Service.ConnService(ServiceName)
			if err == nil {
				break
			}
			fmt.Println(err.Error())
			time.Sleep(1 * time.Second)
		}

		ClientMap[ServiceName] = c
		// 绑定消息队列
		ServiceHandlerMsgMap[ServiceName] = make(map[string]<-chan amqp.Delivery)
		for _, HandlerName := range handler.ServiceMethods[ServiceName] {
			queue, err := Mq.QueueDeclare(HandlerName, false, false, false, false, nil)
			if err != nil {
				panic(err)
			}
			queueMap[HandlerName] = queue
			fmt.Println("bind queue ", queue.Name)
			msg, err := Mq.GetMessages(
				queue.Name, "",
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
					_handler, err := ParseMessage.ParseDelivery(handler.GetHandlersType(ServiceName, HandlerName), delivery)
					if err != nil {
						fmt.Println(err)
						continue
					}
					delivery.Ack(false)
					_handler.GetHandlerName()
					// 调用rpc服务并返回结果
					res, err := c.Handle(ctx, _handler)
					if err != nil {
						fmt.Println(err)
					}
					if res == nil {
						fmt.Println("res is nil")
						continue
					}
					// 发送响应给server
					err = Mq.GivesResponseTo(delivery.ReplyTo, delivery.CorrelationId, res)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}
}
func main() {
	var err error
	err = Settings.Init()
	if err != nil {
		panic(err)
	}
	Init.InitService()
	// 连接消息队列
	err = Mq.NewConnCh(Settings.Config.RabbitMQConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect rabbitmq success at", Settings.Config.RabbitMQConfig.Host, Settings.Config.RabbitMQConfig.Port)
	defer Mq.ConnClose()
	// 启动服务发现
	fmt.Println("start service discovery at", Settings.Config.EtcdConfig.Host, Settings.Config.EtcdConfig.Port)
	StartService(Settings.Config.EtcdConfig.Host, Settings.Config.EtcdConfig.Port)
	// 监听消息队列
	ListenQueues()
	select {}
}
