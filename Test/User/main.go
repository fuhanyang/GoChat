package main

import (
	"Test/UserPb"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rpc/User"
	"time"
)

var servicesName = []string{
	"User",
}

func getServerAddr(svcName string) string {
	s := UserPb.ServiceDiscovery(svcName)
	if s == nil || (s.IP == "" && s.Port == "") {
		return ""
	}
	return s.IP + ":" + s.Port
}
func main() {
	//for _, svcName := range servicesName {
	//	go func() {
	//		err := UserPb.WatchServiceName(svcName)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	//}
	//time.Sleep(1 * time.Second)
	//addr := getServerAddr("User")
	//if addr == "" {
	//	fmt.Println(" service not found")
	//	return
	//}
	//fmt.Println(" service found at ", addr)
	//conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err)
	//}
	//c := UserPb.NewClient(conn)
	//ctx := context.Background()
	//for {
	//	time.Sleep(1 * time.Second)
	//	res, err := c.Register(ctx, &UserPb.RegisterRequest{
	//		Username: "ljq",
	//		Password: "123456",
	//		Ip:       "127.0.0.1",
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(res)
	//	res1, err := c.Login(ctx, &UserPb.LoginRequest{
	//		AccountNum:      res.GetAccountNum(),
	//		Password:        "123456",
	//		PasswordConfirm: "123456",
	//		Ip:              "127.0.0.1",
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(res1)
	//
	//	res2, err := c.LogOff(ctx, &UserPb.LogOffRequest{
	//		AccountNum: res.GetAccountNum(),
	//		Password:   "123456",
	//		Ip:         "127.0.0.1",
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(res2)
	//}

	// 连接消息队列
	err := NewConnCh()
	if err != nil {
		panic(err)
	}
	defer ConnClose()
	//声明回调队列
	for _, ServiceName := range servicesName {
		ServiceHandlerMsgMap[ServiceName] = make(map[string]<-chan amqp.Delivery)
		for _, HandlerName := range HandlersName[ServiceName] {
			queue, err := QueueDeclare("", false, false, false, false, nil)
			if err != nil {
				panic(err)
			}
			queueMap[HandlerName] = queue
			ServiceHandlerMsgMap[ServiceName][HandlerName], err = ch.Consume(
				queue.Name,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			fmt.Println("bind queue ", queue.Name)
		}
	}
	for _, ServiceName := range servicesName {
		for _, HandlerName := range HandlersName[ServiceName] {
			go func() {
				if HandlerName == "Register" {
					for d := range ServiceHandlerMsgMap[ServiceName][HandlerName] {
						fmt.Println("get response :", string(d.Body))
						var rgs User.RegisterResponse
						err := json.Unmarshal(d.Body, &rgs)
						if err != nil {
							fmt.Println("unmarshal err:", err)
							continue
						}
						v1, err := json.Marshal(User.LoginRequest{
							AccountNum:  rgs.AccountNum,
							Password:    "123456",
							Ip:          "127.0.0.1",
							HandlerName: "Login",
						})
						if err != nil {
							fmt.Println("marshal err:", err)
							continue
						}
						err = ch.Publish(
							"",
							"Login",
							false,
							false,
							amqp.Publishing{
								ContentType: "text/plain",
								ReplyTo:     queueMap["Login"].Name,
								Body:        v1,
							},
						)
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
				}
			}()
		}
	}
	for {
		time.Sleep(1 * time.Second)
		v, err := json.Marshal(User.RegisterRequest{
			Username:    "xyx",
			Password:    "123456",
			Ip:          "127.0.0.1",
			HandlerName: "Register",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = ch.Publish(
			"",
			"Register",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				ReplyTo:     queueMap["Register"].Name,
				Body:        v,
			},
		)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}

var conn *amqp.Connection
var ch *amqp.Channel
var ServicesName = []string{
	"User",
}
var HandlersName = map[string][]string{
	"User": {"Register", "Login", "Logoff"},
}
var queueMap = make(map[string]amqp.Queue)
var ServiceHandlerMsgMap = make(map[string]map[string]<-chan amqp.Delivery)

func NewConnCh() error {
	var err error
	conn, err = amqp.Dial("amqp://user:56563096660fc@localhost:8080/")
	if err != nil {
		return err
	}
	ch, err = conn.Channel()
	return err
}
func ConnClose() {
	ch.Close()
	conn.Close()
}
func QueueDeclare(queueName string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
	return q, err
}
